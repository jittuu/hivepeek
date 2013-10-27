package events

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"peek/ds"
	"peek/fb24"
	"time"
)

type EventGoals []*ds.EventGoals

func (eventGoals EventGoals) Find(home, away string) *ds.EventGoals {
	for _, eg := range eventGoals {
		if eg.Home == home && eg.Away == away {
			return eg
		}
	}

	return nil
}

type EventList []*Event

func (events EventList) EventIDMapping(teamMappings map[string]string) EventIDMapping {
	mapping := make(map[string]int64)
	for _, e := range events {
		h := teamMappings[e.Home]
		a := teamMappings[e.Away]
		eventId := fmt.Sprintf("%s-%s-%s", e.StartTime.Format("20060102"), h, a)
		mapping[eventId] = e.Id
	}

	return mapping
}

type EventIDMapping map[string]int64

func (mapping EventIDMapping) Find(home, away string, startTime time.Time) int64 {
	eventId := fmt.Sprintf("%s-%s-%s", startTime.Format("20060102"), home, away)
	return mapping[eventId]
}

type fetchGoalsTask struct {
	context        appengine.Context
	league, season string
}

func (t *fetchGoalsTask) getEvents() (events EventList, err error) {
	dst, keys, err := ds.GetAllEvents(t.context, t.league, t.season)
	if err != nil {
		return
	}
	events = make([]*Event, len(dst))
	for i, e := range dst {
		events[i] = &Event{
			Event: e,
			Id:    keys[i].IntID(),
		}
	}

	return
}

func (t *fetchGoalsTask) getTeamMappings() (mappings map[string]string, err error) {
	dst, _, err := ds.GetAllTeamMappings(t.context)
	if err != nil {
		return
	}

	mappings = make(map[string]string)
	for _, m := range dst {
		mappings[m.Name] = m.MasterName
	}

	return
}

func (t *fetchGoalsTask) exec() error {
	goal_events, err := fb24.Fetch(t.context, t.league, t.season)
	if err != nil {
		return err
	}

	existing_goal_events, err := ds.GetEventGoalsByLeagueAndSeason(t.context, t.league, t.season)
	if err != nil {
		return err
	}
	events, err := t.getEvents()
	if err != nil {
		return err
	}

	teamMappings, err := t.getTeamMappings()
	if err != nil {
		return err
	}

	eventMappings := events.EventIDMapping(teamMappings)

	event_goals := make([]*ds.EventGoals, len(goal_events))
	event_goals_keys := make([]*datastore.Key, len(goal_events))
	for i, ge := range goal_events {
		home := teamMappings[ge.Home()]
		away := teamMappings[ge.Away()]
		if home == "" {
			t.context.Errorf("cannot find mapping for %s", ge.Home())
		}
		if away == "" {
			t.context.Errorf("cannot find mapping for %s", ge.Away())
		}

		eventId := eventMappings.Find(home, away, ge.StartTime())
		existing := EventGoals(existing_goal_events).Find(ge.Home(), ge.Away())
		if existing == nil {
			e := &ds.EventGoals{
				League:    t.league,
				Season:    t.season,
				StartTime: ge.StartTime(),
				Home:      ge.Home(),
				Away:      ge.Away(),
				HomeGoals: ge.HomeGoals(),
				AwayGoals: ge.AwayGoals(),
				EventId:   eventId,
			}

			event_goals[i] = e
			event_goals_keys[i] = datastore.NewIncompleteKey(t.context, "EventGoals", nil)
		} else {
			existing.EventId = eventId
			event_goals[i] = existing
			event_goals_keys[i] = datastore.NewKey(t.context, "EventGoals", "", existing.Id, nil)
		}
	}

	_, err = datastore.PutMulti(t.context, event_goals_keys, event_goals)
	return err
}
