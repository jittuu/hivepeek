package fb24

import (
	"appengine"
	"appengine/urlfetch"
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"time"
)

type Event struct {
	home, away           string
	fmtDate              string
	homeGoals, awayGoals []string
}

func (e *Event) Home() string {
	return e.home
}

func (e *Event) Away() string {
	return e.away
}

func (e *Event) StartTime() time.Time {
	// fb24 use Denmark timezone as base timezone
	loc, _ := time.LoadLocation("Europe/Copenhagen")
	dt, err := time.ParseInLocation("02/01/2006 15:04", e.fmtDate, loc)
	if err != nil {
		panic(err)
	}
	return dt.UTC()
}

func convertGoals(values []string) []int {
	goals := make([]int, len(values))
	for i, g := range values {
		hg := g
		if len(g) > 2 {
			hg = g[:2]
		}
		value, err := strconv.ParseInt(hg, 10, 64)
		if err != nil {
			panic(err)
		}

		goals[i] = int(value)
	}

	return goals
}

func (e *Event) HomeGoals() []int {
	return convertGoals(e.homeGoals)
}

func (e *Event) AwayGoals() []int {
	return convertGoals(e.awayGoals)
}

func getUrl(league, season string) string {
	urlFmt := "http://www.futbol24.com/national/%s/%s/results/?statLR-Page="
	switch league {
	case "epl":
		return fmt.Sprintf(urlFmt, "England/Premier-League", season)
	case "serie-a":
		return fmt.Sprintf(urlFmt, "Italy/Serie-A", season)
	case "bundesliga":
		return fmt.Sprintf(urlFmt, "Germany/Bundesliga", season)
	case "la-liga":
		return fmt.Sprintf(urlFmt, "Spain/Primera-Division", season)
	case "ligue-1":
		return fmt.Sprintf(urlFmt, "France/Ligue-1", season)
	}

	return ""
}

func getWithRetry(c appengine.Context, url string, retry int) (resp *http.Response, err error) {
	for i := 0; i < retry; i++ {
		client := urlfetch.Client(c)
		resp, err = client.Get(url)
		if err == nil && resp.ContentLength > 1000 {
			break
		}

		<-time.After(5 * time.Second)
	}

	return
}

type detailResult struct {
	evt *Event
	err error
}

func fetchDetail(c appengine.Context, url string) <-chan *detailResult {
	ch := make(chan *detailResult, 1)

	go func() {
		resp, err := getWithRetry(c, url, 5)
		if err != nil {
			ch <- &detailResult{nil, err}
			return
		}

		doc, err := gq.NewDocumentFromResponse(resp)
		if err != nil {
			ch <- &detailResult{nil, err}
			return
		}

		event := &Event{}

		event.fmtDate = doc.Find("span.date").Text()

		doc.Find("thead tr td.home, thead tr td.guest").Each(func(i int, s *gq.Selection) {
			if s.Is("td.home") {
				event.home = s.ChildrenFiltered("a").Text()
			} else {
				event.away = s.ChildrenFiltered("a").Text()
			}
		})

		doc.Find("tbody tr.haction1, tbody tr.haction5, tbody tr.gaction1, tbody tr.gaction5, tbody tr.haction4, tbody tr.gaction4").Each(func(i int, s *gq.Selection) {
			h := s.Find("td.home span").Text()
			if h != "" {
				hg := s.ChildrenFiltered("td.status").Text()
				event.homeGoals = append(event.homeGoals, hg)
			} else {
				ag := s.ChildrenFiltered("td.status").Text()
				event.awayGoals = append(event.awayGoals, ag)
			}
		})

		ch <- &detailResult{event, nil}
	}()

	return ch
}

func fetchPage(c appengine.Context, url string) (events []*Event, nextPage string, err error) {
	resp, err := getWithRetry(c, url, 5)
	if err != nil {
		return
	}

	doc, err := gq.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}

	tasks := make([]<-chan *detailResult, 0)
	nextA := doc.Find("div#statLR div.next a.stat_ajax_click")
	nextPageUrl, exist := nextA.Attr("href")
	nextPage = "0"
	if exist {
		nextPage = nextPageUrl[len(nextPageUrl)-1:]
	}

	doc.Find("table.stat tr.status5 a.matchAction").Each(func(i int, s *gq.Selection) {
		link, _ := s.Attr("href")
		link = "http://www.futbol24.com" + link

		tasks = append(tasks, fetchDetail(c, link))
	})

	for _, ch := range tasks {
		result := <-ch
		if result.err == nil {
			evt := result.evt
			events = append(events, evt)
		} else {
			c.Errorf("%s", result.err.Error())
		}
	}

	return
}

func Fetch(c appengine.Context, league, season string) ([]*Event, error) {
	pageNo := "0"
	eventsLen := 0
	result := make([]*Event, 0)
	for {
		url := getUrl(league, season) + pageNo
		c.Infof("fetching url: %s", url)
		events, nextPage, err := fetchPage(c, url)
		if err != nil {
			return nil, err
		}

		for _, e := range events {
			result = append(result, e)
		}
		eventsLen += len(events)
		c.Infof("fetched [%d] events for pageNo: %s", len(events), pageNo)
		pageNo = nextPage
		if pageNo == "0" {
			break
		}
	}
	c.Infof("total fetched events: %d", eventsLen)

	return result, nil
}
