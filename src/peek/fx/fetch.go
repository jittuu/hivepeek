package fx

import (
	"appengine"
	"appengine/urlfetch"
	"bytes"
	"encoding/xml"
	"errors"
	"time"
)

type Sport struct {
	Name    string    `xml:"sport,attr"`
	Leagues []*League `xml:"category"`
}

type League struct {
	Id     int      `xml:"id,attr"`
	Name   string   `xml:"name,attr"`
	Events []*Event `xml:"matches>match"`
}

type Event struct {
	Id      int    `xml:"id,attr"`
	FmtDate string `xml:"formatted_date,attr"`
	FmtTime string `xml:"time,attr"`
	Home    *Team  `xml:"localteam"`
	Away    *Team  `xml:"visitorteam"`
}

type Team struct {
	Id   int    `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

func Fetch(c appengine.Context, league string) (sport *Sport, err error) {
	url, err := getUrl(league)
	if err != nil {
		return
	}

	client := urlfetch.Client(c)
	resp, err := client.Get(url)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return
	}

	sport = &Sport{}
	err = xml.Unmarshal(buffer.Bytes(), sport)
	return
}

func getUrl(league string) (url string, err error) {
	switch league {
	case "epl":
		url = "http://www.goalserve.com/getfeed/421736c766374db393fa4244a760e11a/soccernew/england_shedule"
	case "serie-a":
		url = "http://www.goalserve.com/getfeed/421736c766374db393fa4244a760e11a/soccernew/italy_shedule"
	case "bundesliga":
		url = "http://www.goalserve.com/getfeed/421736c766374db393fa4244a760e11a/soccernew/germany_shedule"
	default:
		err = errors.New("Not supported league")
	}
	return
}

func (e *Event) StartTime() time.Time {
	layout := "02.01.2006 15:04" // 21.09.2013 16:30
	t, _ := time.Parse(layout, e.FmtDate+" "+e.FmtTime)
	return t
}
