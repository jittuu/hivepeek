package events

import (
	"appengine"
	"peek/ds"
	"time"
)

type fixtureQuery struct {
	from    time.Time
	to      time.Time
	context appengine.Context
}

func (f *fixtureQuery) exec() (events []*Event, err error) {
	fixtures, _, err := ds.GetFixtures(f.context, f.from, f.to)
	if err != nil {
		return nil, err
	}

	teamIDs := make([]int64, 0)
	for _, f := range fixtures {
		if f.HomeId != 0 && f.AwayId != 0 {
			teamIDs = append(teamIDs, f.HomeId)
			teamIDs = append(teamIDs, f.AwayId)
		}
	}

	teams, keys, err := ds.GetTeams(f.context, teamIDs)
	if err != nil {
		return nil, err
	}

	teamMaps := make(map[int64]*ds.Team)
	for i, t := range teams {
		teamMaps[keys[i].IntID()] = t
	}

	events = make([]*Event, 0)
	for _, f := range fixtures {
		h := teamMaps[f.HomeId]
		a := teamMaps[f.AwayId]
		var evt *ds.Event
		if h != nil && a != nil {
			evt = &ds.Event{
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
		} else {
			evt = &ds.Event{
				League:    f.League,
				Season:    f.Season,
				StartTime: f.StartTime,
				Home:      f.Home,
				Away:      f.Away,
			}
		}
		events = append(events, &Event{Event: evt})
	}

	startTime := func(e1, e2 *Event) bool {
		return e1.StartTime.Before(e2.StartTime)
	}

	By(startTime).Sort(events)

	return
}
