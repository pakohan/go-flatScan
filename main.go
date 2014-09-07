package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/pakohan/go-libs/flatscan"

	"net/url"

	"appengine"
	"appengine/datastore"
	"appengine/mail"
	_ "appengine/remote_api"
	"appengine/taskqueue"
	"appengine/urlfetch"
)

var emailTemplate *template.Template
var prefTemplate *template.Template

type Zip struct{}

func init() {
	// as long as we never change something in the template, it won't throw an error
	funcMap := template.FuncMap{
		"md5": func(s string) string {
			md5Writer := md5.New()
			io.WriteString(md5Writer, s)
			return fmt.Sprintf("%x", md5Writer.Sum(nil))
		},
	}

	t, err := template.New("email").Funcs(funcMap).Parse(email)
	if err != nil {
		panic(err)
	}
	emailTemplate = t

	pt, err := template.New("pref.html").ParseFiles("pref.html")
	if err != nil {
		panic(err)
	}
	prefTemplate = pt

	http.HandleFunc("/scrape", scrape)
	http.HandleFunc("/listSaved", listSaved)
	http.HandleFunc("/", handle(main))
	http.HandleFunc("/delete", del)
	http.HandleFunc("/worker", checkOffers)
	http.HandleFunc("/index.html", handle(main))
	http.HandleFunc("/pref.html", handle(pref))
}

func sendErrorMail(c appengine.Context, err error) {
	msg := &mail.Message{
		Sender:  "Flat Scan Sender <admin@flat-scan.appspotmail.com>",
		Subject: "Error",
		Body:    err.Error(),
	}

	c.Errorf("%s", msg.Body)
	err = mail.SendToAdmins(c, msg)
}

func errResponse(w http.ResponseWriter, c appengine.Context, err error) {
	w.Header().Add("Content-Type", "application/json")
	http.Error(w, "[]", http.StatusInternalServerError)
	sendErrorMail(c, err)
}

func listSaved(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	w.Header().Add("Content-Type", "application/json")

	dst := make([]flatscan.FlatOffer, 0)
	_, err := datastore.NewQuery(entitiyFlatOffer).GetAll(c, &dst)

	if err != nil {
		errResponse(w, c, err)
		return
	}

	val, err := json.Marshal(dst)
	if err != nil {
		errResponse(w, c, err)
		return
	}

	w.Write(val)
}

func del(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	t := taskqueue.NewPOSTTask("/worker", nil)
	if _, err := taskqueue.Add(c, t, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func checkOffers(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	dst := make([]flatscan.FlatOffer, 0)
	now := time.Now().Unix()
	keys, err := datastore.NewQuery(entitiyFlatOffer).Filter("TimeUpdated <", now-(60*60*24)).GetAll(c, &dst)
	if err != nil {
		c.Errorf("%s", err)
		return
	}

	c.Infof("loaded %d entities to check", len(dst))

	client := urlfetch.Client(c)

	for i, offer := range dst {
		check(keys[i], offer, client, c)
	}
}

func check(key *datastore.Key, offer flatscan.FlatOffer, client *http.Client, c appengine.Context) {
	urlObj, err := url.Parse(fmt.Sprintf("%s%s", base, offer.Url))
	if err != nil {
		c.Errorf(err.Error())
		return
	}

	urlParts := strings.Split(urlObj.Path, "/")
	urlParts = strings.Split(urlParts[len(urlParts)-1], "-")

	link := fmt.Sprintf("%s/anzeigen/s-anzeige/%s", base, urlParts[0])

	resp, err := client.Get(link)
	if err != nil {
		c.Errorf(err.Error())
		return
	}

	_, ok := resp.Request.Header["Referer"]

	if ok {
		c.Infof("Removing Entity with url '%s'", link)
		err := datastore.Delete(c, key)
		if err != nil {
			c.Errorf("%s", err.Error())
		}
	}
}

func AEKey(f flatscan.FlatOffer, con appengine.Context) *datastore.Key {
	return datastore.NewKey(con, "counter", f.Key(), 0, nil)
}
