package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"net/http"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, "/")
		if err != nil {
			http.Error(w, "could not create login url", http.StatusInternalServerError)
			sendErrorMail(c, err)
			return
		}

		http.Redirect(w, r, url, 302)
		return
	}

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
}

func pref(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, "/pref.html")
		if err != nil {
			http.Error(w, "could not create login url", http.StatusInternalServerError)
			sendErrorMail(c, err)
			return
		}

		http.Redirect(w, r, url, 302)
		return
	}

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

	http.ServeFile(w, r, "pref.html")
	return
}
