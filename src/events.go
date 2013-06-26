package peek

import (
	"appengine"
	"appengine/datastore"
	"ds"
	"net/http"
	"sort"
	"time"
)

type Event struct {
	*ds.Event
	Id int64
}

type Team struct {
	*ds.Team
	Id int64
}

func Index(w http.ResponseWriter, r *http.Request) {
	var s = r.FormValue("s")
	if s == "" {
		http.Redirect(w, r, "/events/?s=2012-2013", http.StatusFound)
		return
	}

	c := appengine.NewContext(r)

	if events, err := getEvents(c, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		date := func(e1, e2 *Event) bool {
			return e1.StartTime.Before(e2.StartTime)
		}
		By(date).Sort(events)
		renderTemplate(w, events, "templates/events.html")
	}
}

func Calc(w http.ResponseWriter, r *http.Request) {
	var s = r.FormValue("s")

	c := appengine.NewContext(r)

	if events, err := getEvents(c, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		batches := make(map[time.Time][]*Event, 0)
		for _, e := range events {
			batches[e.StartTime] = append(batches[e.StartTime], e)
		}

		teams, err := getTeams(c)
		findTeam := func(n string) *Team {
			for _, t := range teams {
				if t.Name == n {
					return t
				}
			}

			return nil
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dates := sortedDates(batches)
		for _, d := range dates {
			for _, e := range batches[d] {
				h := findTeam(e.Home)
				a := findTeam(e.Away)
				var hr, ar *ds.SeasonRating
				for i, _ := range h.Ratings {
					if h.Ratings[i].Season == s {
						hr = &(h.Ratings[i])
					}
				}
				for i, _ := range a.Ratings {
					if a.Ratings[i].Season == s {
						ar = &(a.Ratings[i])
					}
				}

				e.HRating = hr.Overall
				e.ARating = ar.Overall

				switch {
				case e.HGoal > e.AGoal:
					transfer(hr, ar, 10)
				case e.HGoal == e.AGoal:
					transfer(ar, hr, 5)
				case e.HGoal < e.AGoal:
					transfer(ar, hr, 20)
				}

				datastore.Put(c, datastore.NewKey(c, "Event", "", e.Id, nil), e.Event)
				datastore.Put(c, datastore.NewKey(c, "Team", "", h.Id, nil), h.Team)
				datastore.Put(c, datastore.NewKey(c, "Team", "", a.Id, nil), a.Team)
			}
		}

		http.Redirect(w, r, "/events/?s="+s, http.StatusFound)
	}
}

func transfer(w *ds.SeasonRating, l *ds.SeasonRating, percent int) {
	amt := l.Overall * percent / 100
	w.Overall += amt
	l.Overall -= amt
	return
}

func getEvents(c appengine.Context, s string) ([]*Event, error) {
	if dst, keys, err := ds.GetAllEvents(c, s); err != nil {
		return nil, err
	} else {
		events := make([]*Event, len(dst))
		for i, e := range dst {
			events[i] = &Event{
				Event: e,
				Id:    keys[i].IntID()}
		}
		return events, nil
	}
}

func getTeams(c appengine.Context) ([]*Team, error) {
	if dst, keys, err := ds.GetAllTeams(c); err != nil {
		return nil, err
	} else {
		teams := make([]*Team, len(dst))
		for i, t := range dst {
			teams[i] = &Team{
				Team: t,
				Id:   keys[i].IntID()}
		}
		return teams, nil
	}
}

func sortedDates(m map[time.Time][]*Event) []time.Time {
	dates := make(byTime, len(m))
	i := 0
	for k, _ := range m {
		dates[i] = k
		i++
	}

	sort.Sort(dates)

	return dates
}

type byTime []time.Time

func (s byTime) Less(i, j int) bool { return s[i].Before(s[j]) }
func (s byTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byTime) Len() int           { return len(s) }

