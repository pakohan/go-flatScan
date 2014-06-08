package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"bytes"
	"net/http"
	"strings"
	"time"
)

type handler func(w http.ResponseWriter, r *http.Request)
type handlerInfo struct {
	loginURL string
	loggedIn handler
}

var pref = handlerInfo{
	"/pref.html",
	func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		var pref *Setting
		key := datastore.NewKey(c, userEntitiy, u.Email, 0, nil)

		c.Infof(r.Method)

		if strings.EqualFold(r.Method, "POST") {
			r.ParseForm()
			pref = NewSetting(r.Form, u.Email)
			c.Infof("%+v", *pref)
			_, err := datastore.Put(c, key, pref)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		} else {
			pref = &Setting{}
			err := datastore.Get(c, key, pref)
			if err != nil && err != datastore.ErrNoSuchEntity {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		buf := bytes.NewBufferString("")
		err := prefTemplate.Execute(buf, pref)
		if err != nil {
			c.Errorf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(buf.Bytes()))
		return
	},
}

var main handlerInfo = handlerInfo{
	"/main.html",
	func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		setting := Setting{}
		key := datastore.NewKey(c, userEntitiy, u.Email, 0, nil)
		err := datastore.Get(c, key, &setting)

		if err == datastore.ErrNoSuchEntity {
			_, err = datastore.Put(c, key, &setting)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		} else if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.ServeFile(w, r, "index.html")
		return
	},
}

func handle(h handlerInfo) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		if u == nil {
			url, err := user.LoginURL(c, h.loginURL)
			if err != nil {
				http.Error(w, "could not create login url", http.StatusInternalServerError)
				sendErrorMail(c, err)
				return
			}

			http.Redirect(w, r, url, 302)
			return
		} else {
			h.loggedIn(w, r)
		}
	}
}
