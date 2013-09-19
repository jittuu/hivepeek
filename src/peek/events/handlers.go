package events

import (
	"appengine"
	"appengine/memcache"
	"appengine/user"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"peek"
	"peek/ds"
	"strconv"
	"time"
)

func newUpload(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/upload.html")
}

func upload(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := parseEvents(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	season := r.FormValue("season")
	league := r.FormValue("league")
	update, _ := strconv.ParseBool(r.FormValue("update"))

	t := &uploadTask{
		context: c,
		events:  events,
		season:  season,
		league:  league,
		update:  update,
	}
	if err := t.exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := memcache.Flush(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := "/events/" + league + "?s=" + season
	http.Redirect(w, r, url, http.StatusFound)
}

const layout = "2006-01-02"

func index(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	url := "/events/epl?s=" + getSeason(now) + "&d=" + now.Format(layout)
	http.Redirect(w, r, url, http.StatusFound)
	return
}

func fixture(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	y, m, d := time.Now().Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, time.Now().Location())
	next_3days := today.Add(3 * 24 * time.Hour)
	fixtures, _, err := ds.GetFixtures(c, today, next_3days)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teamIDs := make([]int64, 0)
	for _, f := range fixtures {
		if f.HomeId != 0 && f.AwayId != 0 {
			teamIDs = append(teamIDs, f.HomeId)
			teamIDs = append(teamIDs, f.AwayId)
		}
	}

	teams, keys, err := ds.GetTeams(c, teamIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	teamMaps := make(map[int64]*ds.Team)
	for i, t := range teams {
		teamMaps[keys[i].IntID()] = t
	}

	events := make([]*Event, 0)
	for _, f := range fixtures {
		h := teamMaps[f.HomeId]
		a := teamMaps[f.AwayId]
		if h != nil || a != nil {
			evt := &ds.Event{
				League:         f.League,
				Season:         f.Season,
				StartTime:      f.StartTime,
				Home:           h.Name,
				Away:           a.Name,
				HRating:        h.OverallRating,
				HRatingLen:     h.OverallRatingLen,
				HNetRating:     h.HomeNetRating,
				HNetRatingLen:  h.HomeNetRatingLen,
				HFormRating:    h.FormRating(),
				HFormRatingLen: len(h.LastFiveMatchRating),
				ARating:        a.OverallRating,
				ARatingLen:     a.OverallRatingLen,
				ANetRating:     a.AwayNetRating,
				ANetRatingLen:  a.AwayNetRatingLen,
				AFormRating:    a.FormRating(),
				AFormRatingLen: len(a.LastFiveMatchRating),
			}
			events = append(events, &Event{Event: evt})
		}
	}

	peek.RenderTemplate(w, events, "templates/fixtures.html")
}

func league(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := vars["league"]
	season := r.FormValue("s")
	d := r.FormValue("d")

	if season == "" || d == "" {
		now := time.Now()
		if season == "" {
			season = getSeason(now)
		}
		url := "/events/" + league + "?s=" + season + "&d=" + now.Format(layout)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

	date, _ := time.Parse(layout, d)
	events, _ := getEventsByWeek(c, league, season, date)

	gw := &GameWeek{
		Events:      events,
		PreviousUrl: "/events/" + league + "?s=" + season + "&d=" + date.AddDate(0, 0, -7).Format(layout),
		NextUrl:     "/events/" + league + "?s=" + season + "&d=" + date.AddDate(0, 0, 7).Format(layout),
		League:      league,
		Season:      season,
		Date:        date,
		IsAdmin:     user.IsAdmin(c),
	}

	peek.RenderTemplate(w, gw, "templates/events.html")
	return
}

func fetchView(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/fetch.html")
}

func fetch(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	league := r.FormValue("league")
	task := &fetchTask{
		context: c,
		league:  league,
	}
	err := task.exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := "/events/fixture"
	http.Redirect(w, r, url, http.StatusFound)
	return
}

func calcView(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/calc.html")
}

func calc(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	league := r.FormValue("league")
	season := r.FormValue("season")

	t := &calcTask{
		context: c,
		season:  season,
		league:  league,
	}

	if err := t.exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := memcache.Flush(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := "/events/" + league + "?s=" + season
	http.Redirect(w, r, url, http.StatusFound)
}

func resetView(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/reset.html")
}

func reset(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	league := r.FormValue("league")
	season := r.FormValue("season")
	t := &resetTask{
		context: c,
		season:  season,
		league:  league,
	}

	if err := t.exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := memcache.Flush(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := "/events/" + league + "?s=" + season
	http.Redirect(w, r, url, http.StatusFound)
}

func runView(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/run.html")
}

func run(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	league := r.FormValue("league")
	season := r.FormValue("season")
	diff, _ := strconv.ParseFloat(r.FormValue("diff"), 64)
	min, _ := strconv.ParseFloat(r.FormValue("min"), 64)
	max, _ := strconv.ParseFloat(r.FormValue("max"), 64)

	t := &runTask{
		context:  c,
		w:        ioutil.Discard,
		season:   season,
		league:   league,
		diff:     diff,
		minPrice: min,
		maxPrice: max,
	}

	result, err := t.exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vm := &RunTaskResult{Results: result}

	vm.Bets, vm.Profit = result.Profit()
	peek.RenderTemplate(w, vm, "templates/runresult.html")
}

func getEventsByWeek(c appengine.Context, league, season string, date time.Time) ([]*Event, error) {
	key := fmt.Sprintf("%s-%s-%s", league, season, date.Format(layout))
	cached, err := memcache.Get(c, key)
	if err != nil && err != memcache.ErrCacheMiss {
		return nil, err
	}

	if err == nil {
		// cache hit
		var b bytes.Buffer
		b.Write(cached.Value)
		dec := gob.NewDecoder(&b)
		events := make([]*Event, 0)
		if err = dec.Decode(&events); err != nil {
			return nil, err
		}

		return events, nil
	} else {
		// cache miss
		start, end := weekRange(date)
		dst, keys, err := ds.GetAllEventsByDateRange(c, league, season, start, end)

		if err != nil {
			return nil, err
		}

		events := make([]*Event, len(dst))
		for i, e := range dst {
			events[i] = &Event{
				Event: e,
				Id:    keys[i].IntID(),
			}
		}
		if err = setEventsToCache(c, key, events); err != nil {
			return nil, err
		}

		return events, nil
	}
}

func setEventsToCache(c appengine.Context, key string, events []*Event) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(events)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:   key,
		Value: b.Bytes(),
	}
	if err = memcache.Set(c, item); err != nil {
		return err
	}

	return nil
}

func weekRange(date time.Time) (start, end time.Time) {
	y, m, d := date.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, date.Location())
	days := (-1 * (int(date.Weekday()) + 1)) % 7
	start = today.AddDate(0, 0, days)
	end = start.AddDate(0, 0, 6)
	return
}

func getSeason(date time.Time) string {
	y, m, _ := date.Date()
	start, end := y, y
	if m > 7 {
		end++
	} else {
		start--
	}

	return fmt.Sprintf("%d-%d", start, end)
}
