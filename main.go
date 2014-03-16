package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/mail"
	"appengine/urlfetch"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pakohan/go-libs/flatscan"
	"net/http"
	"strconv"
	"text/template"
)

const (
	base       string = "http://kleinanzeigen.ebay.de"
	searchSite string = "%s/anzeigen/s-wohnung-mieten/berlin/anzeige:angebote/seite:%d/c203l3331"
)

var emailTemplate *template.Template

const email string = `
You have a new offer:
{{if gt .RentN 0}}Rent: {{.RentN}}
{{end}}Adresse: {{if gt (len .Street) 0}}{{.Street}}
{{end}}{{if gt (len .District) 0}}{{.District}}
{{end}}{{if gt .Zip 0}}{{.Zip}} {{end}}Berlin
Rooms: {{.Rooms}}
Size: {{.Size}}
Url: http://kleinanzeigen.ebay.de/{{.Url}}

Description: {{.Description}}`

func init() {
	tmpl, err := template.New("email").Parse(email)
	if err != nil {
		panic(err)
	}
	emailTemplate = tmpl

	http.HandleFunc("/scrape", scrape)
	http.HandleFunc("/initialScrape", initialScrape)
	http.HandleFunc("/listSaved", listSaved)
}

func scrape(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	searchUrl := fmt.Sprintf(searchSite, base, 1)
	loadList(searchUrl, c)
}

func listSaved(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	dst := make([]flatscan.FlatOffer, 0)

	_, err := datastore.NewQuery("FlatOffer").Filter("Valid =", true).GetAll(c, &dst)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	val, err := json.Marshal(dst)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(val)
}

func initialScrape(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	pages, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)
	if err != nil {
		panic(err.Error())
		return
	}

	for i := 1; i < int(pages); i++ {
		searchUrl := fmt.Sprintf(searchSite, base, i)
		loadList(searchUrl, c)
	}
}

func loadList(url string, c appengine.Context) {
	client := urlfetch.Client(c)
	doc, err := LoadDocumentGAE(url, client)
	if err != nil {
		panic(err.Error())
	}

	for _, path := range flatscan.ExtractLinks(doc) {
		key := datastore.NewKey(c, "FlatOffer", path, 0, nil)
		amount, err := CheckAmountGAE(key, c)
		if amount > 0 || err != nil {
			continue
		}

		offerUrl := fmt.Sprintf("%s%s", base, path)
		doc, err = LoadDocumentGAE(offerUrl, client)
		if err != nil {
			panic(err)
		}

		offer := flatscan.GetOffer(doc)
		offer.Url = path

		offer.Valid = flatscan.CheckOffer(offer)
		if offer.Valid {
			buf := bytes.NewBufferString("")
			err = emailTemplate.Execute(buf, offer)
			if err != nil {
				panic(err)
			}

			msg := &mail.Message{
				Sender:  "Flat Scan Sender <patrick.kohan@gmail.com>",
				Subject: "Found a Flat",
				Body:    buf.String(),
			}

			err = mail.SendToAdmins(c, msg)
			if err != nil {
				c.Errorf("Couldn't send email: %v", err)
			}
		}

		key, err = datastore.Put(c, key, offer)
	}
}

func CheckAmountGAE(key *datastore.Key, c appengine.Context) (amount int, err error) {
	return datastore.NewQuery("FlatOffer").Filter("__key__ =", key).Count(c)
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
