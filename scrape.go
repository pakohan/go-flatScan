package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/pakohan/go-libs/flatscan"

	"appengine"
	"appengine/datastore"
	"appengine/mail"
	"appengine/urlfetch"
)

func scrape(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	i := 0
	j := 0
	counter := 1
	s, err := GetSettings(c)
	if err != nil {
		c.Errorf(err.Error())
		return
	}
	for i == 0 && j < 5 {
		j++
		searchUrl := fmt.Sprintf(searchSite, base, counter)
		counter++
		var err error
		i, err = loadList(searchUrl, c, s)
		if err != nil {
			sendErrorMail(c, err)
			return
		}
	}
}

func loadList(url string, c appengine.Context, s []Setting) (i int, err error) {
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
			i++
			continue
		}

		offerUrl := fmt.Sprintf("%s%s", base, offerPath)
		doc, err = LoadDocumentGAE(offerUrl, client)
		if err != nil {
			c.Errorf(err.Error())
			i++
			continue
		}

		offer, err := flatscan.GetOffer(doc, c)
		if err != nil {
			c.Errorf("%s: %s", err.Error(), offerPath)
			i++
			continue
		}

		offer.Url = offerPath
		c.Infof("Title: '%s', url: %s", offer.Title, offer.Url)

		offer.ID = md5Sum

		for _, setting := range s {
			b, err := setting.CheckOffer(*offer)
			if b {
				buf := bytes.NewBufferString("")
				err = emailTemplate.Execute(buf, offer)
				if err != nil {
					c.Errorf(err.Error())
				}

				msg := &mail.Message{
					Sender:  "Flat Scan Sender <admin@flat-scraper.appspotmail.com>",
					To:      []string{setting.Email},
					Subject: "Found a Flat",
					Body:    buf.String(),
				}

				c.Infof(msg.Body)

				err = mail.Send(c, msg)
				if err != nil {
					c.Errorf(err.Error())
					sendErrorMail(c, err)
					//return i, err
				} else {
					c.Infof("Mail sent %+v", *msg)
				}
			}
		}

		key, err = datastore.Put(c, key, offer)
		if err != nil {
			c.Errorf(err.Error())
			sendErrorMail(c, err)
			//return i, err
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
