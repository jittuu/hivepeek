package events

import (
	"appengine"
	"appengine/memcache"
	"appengine/taskqueue"
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

func pull(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	season := r.FormValue("season")
	league := r.FormValue("league")
	update, _ := strconv.ParseBool(r.FormValue("update"))

	DelayPull.Call(c, league, season, update)

	http.Redirect(w, r, "/events/qstats", http.StatusFound)
	return
}

func getPull(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/pull.html")
}

func qstats(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	stats, err := taskqueue.QueueStats(c, []string{"default"}, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := &QueueStats{}
	for _, s := range stats {
		data.Total += s.Tasks
		data.Running += s.InFlight
		data.JustFinished += s.Executed1Minute
	}
	peek.RenderTemplate(w, data, "templates/qstats.html")
}

const layout = "2006-01-02"

func index(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	url := "/events/epl?s=" + getSeason(now) + "&d=" + now.Format(layout)
	http.Redirect(w, r, url, http.StatusFound)
	return
}

func today(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	tmr := time.Now().AddDate(0, 0, 1)
	q := &fixtureQuery{
		from:    time.Now(),
		to:      tmr,
		context: c,
	}

	events, err := q.exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	peek.RenderTemplate(w, events, "templates/fixtures.html")
}

func fixture(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	next_month := time.Now().AddDate(0, 1, 0)
	q := &fixtureQuery{
		from:    time.Now(),
		to:      next_month,
		context: c,
	}

	events, err := q.exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	peek.RenderTemplate(w, events, "templates/fixtures.html")
}

func download(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	league := vars["league"]
	season := vars["season"]

	var buf bytes.Buffer

	dl := &dlTask{
		context: c,
		league:  league,
		season:  season,
		w:       &buf,
	}

	if err := dl.exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fn := fmt.Sprintf("%s-%s.csv", league, season)
	w.Header().Add("Content-type", "text/csv")
	w.Header().Add("Content-disposition", "attachment;filename="+fn)
	w.Write(buf.Bytes())

	return
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

	DelayFetch.Call(c, league)

	http.Redirect(w, r, "/events/qstats", http.StatusFound)
	return
}

func calcView(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	peek.RenderTemplate(w, nil, "templates/calc.html")
}

func calc(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	league := r.FormValue("league")
	season := r.FormValue("season")

	DelayCalc.Call(c, league, season)

	http.Redirect(w, r, "/events/qstats", http.StatusFound)
	return
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
