package events

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"peek/ds"
)

type calcTask struct {
	context appengine.Context
	season  string
	league  string
}

type ErrMissingGoalDetails struct {
	Home, Away string
}

func (e *ErrMissingGoalDetails) Error() string {
	return fmt.Sprintf("unmatch goal info for event (%s vs %s)", e.Home, e.Away)
}

func (t *calcTask) getAllEvents() ([]*Event, error) {
	dst, keys, err := ds.GetAllEvents(t.context, t.league, t.season)
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

	startTime := func(e1, e2 *Event) bool {
		return e1.StartTime.Before(e2.StartTime)
	}

	By(startTime).Sort(events)

	return events, nil
}

func (t *calcTask) exec() error {
	events, err := t.getAllEvents()

	if err != nil {
		return err
	}

	for _, e := range events {
		err = t.execEvent(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *calcTask) execEvent(e *Event) error {
	if e.Calculated() {
		return nil
	}

	if e.HGoal != len(e.HGoals) {
		return &ErrMissingGoalDetails{e.Home, e.Away}
	}
	if e.AGoal != len(e.AGoals) {
		return &ErrMissingGoalDetails{e.Home, e.Away}
	}

	h, a, err := t.getTeams(e)
	if err != nil {
		return err
	}

	e.HRating = h.OverallRating
	e.HRatingLen = h.OverallRatingLen
	e.HNetRating = h.HomeNetRating
	e.HNetRatingLen = h.HomeNetRatingLen
	e.HFormRating = h.FormRating()
	e.HFormRatingLen = len(h.LastFiveMatchRating)
	e.HGoalsFor = h.GoalsFor
	e.HGoalsAgainst = h.GoalsAgainst
	e.ARating = a.OverallRating
	e.ARatingLen = a.OverallRatingLen
	e.ANetRating = a.AwayNetRating
	e.ANetRatingLen = a.AwayNetRatingLen
	e.AFormRating = a.FormRating()
	e.AFormRatingLen = len(a.LastFiveMatchRating)
	e.AGoalsFor = a.GoalsFor
	e.AGoalsAgainst = a.GoalsAgainst

	switch {
	case e.HGoal > e.AGoal:
		amt := t.transfer(h, a, 0.1)
		h.HomeNetRating += amt
		a.AwayNetRating -= amt
	case e.HGoal == e.AGoal:
		amt := t.transfer(a, h, 0.05)
		h.HomeNetRating -= amt
		a.AwayNetRating += amt
	case e.HGoal < e.AGoal:
		amt := t.transfer(a, h, 0.2)
		h.HomeNetRating -= amt
		a.AwayNetRating += amt
	}

	h.OverallRatingLen += 1
	h.HomeNetRatingLen += 1
	a.OverallRatingLen += 1
	a.AwayNetRatingLen += 1

	h.CalcEventGoalsSector(e.HGoals, e.AGoals)
	a.CalcEventGoalsSector(e.AGoals, e.HGoals)

	ch := make(chan bool)
	go func() {
		datastore.Put(t.context, datastore.NewKey(t.context, "Event", "", e.Id, nil), e.Event)
		ch <- true
	}()
	go func() {
		datastore.Put(t.context, datastore.NewKey(t.context, "Team", "", h.Id, nil), h.Team)
		ch <- true
	}()
	go func() {
		datastore.Put(t.context, datastore.NewKey(t.context, "Team", "", a.Id, nil), a.Team)
		ch <- true
	}()

	<-ch
	<-ch
	<-ch
	return nil
}

func (t *calcTask) transfer(w *Team, l *Team, percent float64) float64 {
	amt := l.OverallRating * percent
	w.OverallRating += amt
	l.OverallRating -= amt
	w.AddRating(amt)
	l.AddRating(amt * -1)
	return amt
}

func (t *calcTask) getTeams(e *Event) (h, a *Team, err error) {
	ch := make(chan error)

	dsH := &ds.Team{}
	go func() {
		errH := datastore.Get(
			t.context,
			datastore.NewKey(t.context, "Team", "", e.HomeId, nil),
			dsH)
		ch <- errH
	}()

	dsA := &ds.Team{}
	go func() {
		errA := datastore.Get(
			t.context,
			datastore.NewKey(t.context, "Team", "", e.AwayId, nil),
			dsA)
		ch <- errA
	}()

	err = <-ch
	if err2 := <-ch; err2 != nil {
		err = err2
	}

	return &Team{Team: dsH, Id: e.HomeId}, &Team{Team: dsA, Id: e.AwayId}, err
}
