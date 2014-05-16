package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/mail"
	"appengine/urlfetch"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pakohan/go-libs/flatscan"
	"io"
	"net/http"
	"strconv"
	"text/template"
)

const (
	base             string = "http://kleinanzeigen.ebay.de"
	searchSite       string = "%s/anzeigen/s-wohnung-mieten/berlin/anzeige:angebote/seite:%d/c203l3331"
	entitiyFlatOffer string = "FlatOffer"
	zipEntity        string = "ZIP"
	userEntitiy      string = "USER"
	email            string = `
You have a new offer:
{{if gt .RentN 0.0}}Rent: {{.RentN}}
{{end}}Adresse: {{if gt (len .Street) 0}}{{.Street}}
{{end}}{{if gt (len .District) 0}}{{.District}}
{{end}}{{if gt .Zip 0}}{{.Zip}} {{end}}Berlin
Rooms: {{.Rooms}}
Size: {{.Size}}
Url: http://kleinanzeigen.ebay.de{{.Url}}

Remove ZIP: http://flat-scan.appspot.com/removeZip?ID={{.Zip}}
Set as invalid: http://flat-scan.appspot.com/toggleOffer?ID={{md5 .Url}}&valid=false

Description: {{.Description}}`
)

var emailTemplate *template.Template

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
	http.HandleFunc("/toggleOffer", toggleOffer)
	http.HandleFunc("/removeZip", removeZip)
	http.HandleFunc("/", index)
	http.HandleFunc("/index.html", index)
	http.HandleFunc("/pref.html", pref)
}

func scrape(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	i := 27
	counter := 1
	for i == 27 {
		searchUrl := fmt.Sprintf(searchSite, base, counter)
		counter++
		var err error
		_, err = loadList(searchUrl, c)
		if err != nil {
			sendErrorMail(c, err)
			return
		}
	}
}

type Zip struct{}

func sendErrorMail(c appengine.Context, err error) {
	msg := &mail.Message{
		Sender:  "Flat Scan Sender <admin@flat-scan.appspotmail.com>",
		Subject: "Error",
		Body:    err.Error(),
	}

	c.Errorf("%s", err)
	err = mail.SendToAdmins(c, msg)
}

func listSaved(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	dst := make([]flatscan.FlatOffer, 0)
	id, err := strconv.ParseInt(r.FormValue("scope"), 10, 64)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, "[]", http.StatusInternalServerError)
		sendErrorMail(c, err)
		return
	}

	switch id {
	case 0:
		_, err = datastore.NewQuery(entitiyFlatOffer).GetAll(c, &dst)
	case 1:
		_, err = datastore.NewQuery(entitiyFlatOffer).Filter("Valid =", true).GetAll(c, &dst)
	case 2:
		_, err = datastore.NewQuery(entitiyFlatOffer).Filter("Valid =", false).GetAll(c, &dst)
	}

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, "[]", http.StatusInternalServerError)
		sendErrorMail(c, err)
		return
	}

	keys, err := datastore.NewQuery(zipEntity).KeysOnly().GetAll(c, nil)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, "[]", http.StatusInternalServerError)
		sendErrorMail(c, err)
		return
	}

	zipMap := make(map[int64]bool)
	for _, v := range keys {
		zipMap[v.IntID()] = true
	}

	dst2 := make([]flatscan.FlatOffer, 0)
	for _, v := range dst {
		if !zipMap[v.Zip] {
			dst2 = append(dst2, v)
		}
	}

	val, err := json.Marshal(dst2)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, "[]", http.StatusInternalServerError)
		sendErrorMail(c, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(val)
}

func toggleOffer(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := r.FormValue("ID")
	valid := r.FormValue("valid")
	isValid, err := strconv.ParseBool(valid)
	if err != nil {
		http.Error(w, fmt.Sprintf("parameter valid '%s' not valid", valid), http.StatusBadRequest)
		sendErrorMail(c, err)
		return
	}

	key := datastore.NewKey(c, entitiyFlatOffer, id, 0, nil)
	dst := make([]flatscan.FlatOffer, 0)

	// change: get all offers and update the one to change, return the rest
	_, err = datastore.NewQuery(entitiyFlatOffer).Filter("__key__ =", key).GetAll(c, &dst)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't load offer"), http.StatusInternalServerError)
		sendErrorMail(c, err)
		return
	}

	if len(dst) > 0 {
		offer := dst[0]
		offer.Valid = isValid
		key, err = datastore.Put(c, key, &offer)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't save offer"), http.StatusInternalServerError)
			sendErrorMail(c, err)
		}
	}
}

func removeZip(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	zipString := r.FormValue("ID")
	zip, err := strconv.ParseInt(zipString, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("parameter valid '%s' not valid", zipString), http.StatusBadRequest)
		sendErrorMail(c, err)
		return
	}

	key := datastore.NewKey(c, zipEntity, "", zip, nil)
	key, err = datastore.Put(c, key, &Zip{})
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't blacklist zip '&s'", zipString), http.StatusInternalServerError)
		sendErrorMail(c, err)
	}
}

func loadList(url string, c appengine.Context) (i int, err error) {
	s, err := GetSettings(c)
	i = 0
	client := urlfetch.Client(c)
	var doc *goquery.Document
	doc, err = LoadDocumentGAE(url, client)
	if err != nil {
		return
	}

	for _, offerPath := range flatscan.ExtractLinks(doc) {
		md5Writer := md5.New()
		io.WriteString(md5Writer, offerPath)
		md5Sum := fmt.Sprintf("%x", md5Writer.Sum(nil))

		key := datastore.NewKey(c, entitiyFlatOffer, md5Sum, 0, nil)

		amount, err := CheckAmountGAE(key, c)

		if amount > 0 || err != nil {
			continue
		} else {
			i++
		}

		offerUrl := fmt.Sprintf("%s%s", base, offerPath)
		doc, err = LoadDocumentGAE(offerUrl, client)
		if err != nil {
			return i, err
		}

		offer, err := flatscan.GetOffer(doc, c)
		if err != nil {
			return i, err
		}

		offer.Url = offerPath
		offer.ID = md5Sum

		for _, setting := range s {
			if setting.CheckOffer(*offer) {
				buf := bytes.NewBufferString("")
				err = emailTemplate.Execute(buf, offer)
				if err != nil {
					return i, err
				}

				msg := &mail.Message{
					Sender:  "Flat Scan Sender <admin@flat-scan.appspotmail.com>",
					To:      []string{setting.Email},
					Subject: "Found a Flat",
					Body:    buf.String(),
				}

				c.Infof(msg.Body)

				err = mail.Send(c, msg)
				if err != nil {
					return i, err
				}
			}
		}

		key, err = datastore.Put(c, key, offer)
		if err != nil {
			c.Infof(fmt.Sprintf("%+v", offer))
			return i, err
		}
	}

	return i, nil
}

func CheckAmountGAE(key *datastore.Key, c appengine.Context) (amount int, err error) {
	return datastore.NewQuery(entitiyFlatOffer).Filter("__key__ =", key).Count(c)
}

func LoadDocumentGAE(url string, client *http.Client) (doc *goquery.Document, err error) {
	resp, err := client.Get(url)
	if err != nil {
		doc = nil
		return
	}

	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		doc = nil
	}

	return
}

/*
func LoadDocument(url string) (doc *goquery.Document, err error) {
	resp, err := http.Get(url)
	if err != nil {
		doc = nil
		return
	}

	doc, err = goquery.NewDocumentFromResponse(resp)
	if err != nil {
		doc = nil
	}

	return
}

func CheckAmount(path string) (amount int, err error) {
	return 0, nil
}
*/
