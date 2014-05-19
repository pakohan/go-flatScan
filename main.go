package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/mail"
	"appengine/urlfetch"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pakohan/go-libs/flatscan"
	"io"
	"net/http"
	"text/template"
	"time"
)

var emailTemplate *template.Template

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

	http.HandleFunc("/scrape", scrape)
	http.HandleFunc("/listSaved", listSaved)
	http.HandleFunc("/", index)
	http.HandleFunc("/index.html", index)
	http.HandleFunc("/pref.html", pref)
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

func checkOffers(c appengine.Context) {
	dst := make([]flatscan.FlatOffer, 0)
	now := time.Now().Unix()
	keys, err := datastore.NewQuery(entitiyFlatOffer).Filter("TimeUpdated <", now-(60*60*24)).GetAll(c, &dst)
	if err != nil {
		c.Errorf("%s", err)
		return
	}

	client := urlfetch.Client(c)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("DO NOT FOLLOW REDIRECT")
	}

	for i, offer := range dst {
		resp, _ := client.Get(fmt.Sprintf("%s%s", base, offer.Url))

		if resp.StatusCode == 301 {
			c.Infof("Removing Entity with url '%s'", offer.Url)
			err = datastore.Delete(c, keys[i])
			if err != nil {
				c.Errorf("%s", err.Error())
			}
		}
	}
}
