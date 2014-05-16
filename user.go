package main

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"fmt"
	"github.com/pakohan/go-libs/flatscan"
	"net/url"
	"strconv"
	"strings"
)

type Setting struct {
	Districts string
	MinPrice  int64
	MaxPrice  int64
	MinRooms  float64
	MaxRooms  float64
	Email     string
}

func getRangeValues(formValue string) (beginRange, endRange string, err error) {
	if !strings.Contains(formValue, ";") {
		err = errors.New(fmt.Sprintf("wrong entry '%+v'!", formValue))
		return
	}

	parts := strings.Split(formValue, ";")
	beginRange = parts[0]
	endRange = parts[1]
	return
}

func NewSetting(form url.Values, mail string) (pref *Setting) {
	pref = &Setting{Email: mail}
	pref.ChangeSetting(form)

	return
}

func (s *Setting) ChangeSetting(form url.Values) (numErrors int) {
	s.Districts = strings.ToLower(strings.Join(form["district"], ";"))

	beginRange, endRange, err := getRangeValues(form.Get("price"))
	if err == nil {
		s.MinPrice, err = strconv.ParseInt(beginRange, 10, 64)
		if err != nil {
			numErrors++
		}

		s.MaxPrice, err = strconv.ParseInt(endRange, 10, 64)
		if err != nil {
			numErrors++
		}
	} else {
		numErrors++
	}

	beginRange, endRange, err = getRangeValues(form.Get("rooms"))
	if err == nil {
		s.MinRooms, err = strconv.ParseFloat(beginRange, 64)
		if err != nil {
			numErrors++
		}

		s.MaxRooms, err = strconv.ParseFloat(endRange, 64)
		if err != nil {
			numErrors++
		}
	} else {
		numErrors++
	}

	return
}

func (s Setting) CheckOffer(offer flatscan.FlatOffer) (interested bool) {
	interested = offer.Rooms > 0 && (offer.Rooms >= s.MinRooms) && (offer.Rooms <= s.MaxRooms)
	if !interested {
		return
	}

	interested = offer.RentN > 0 && (int64(offer.RentN) >= s.MinPrice) && (int64(offer.RentN) <= s.MaxPrice)
	if !interested {
		return
	}

	interested = len(offer.District) > 0 && strings.Contains(s.Districts, strings.ToLower(offer.District))
	return
}

func GetSettings(c appengine.Context) (settings []Setting, err error) {
	_, err = datastore.NewQuery(userEntitiy).GetAll(c, &settings)
	return
}
