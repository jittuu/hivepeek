package events

import (
	"appengine"
	"appengine/datastore"
	"peek/ds"
)

type calcTask struct {
	context appengine.Context
	events  []*Event
	season  string
	league  string
}

func (t *calcTask) exec() error {
	for _, e := range t.events {
		if e.HRating > 0 || e.ARating > 0 {
			continue
		}

		h, a, err := t.getTeams(e)
		if err != nil {
			return err
		}

		e.HRating = h.OverallRating
		e.HNetRating = h.HomeNetRating
		e.HFormRating = h.FormRating()
		e.ARating = a.OverallRating
		e.ANetRating = a.AwayNetRating
		e.AFormRating = a.FormRating()

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

		datastore.Put(t.context, datastore.NewKey(t.context, "Event", "", e.Id, nil), e.Event)
		datastore.Put(t.context, datastore.NewKey(t.context, "Team", "", h.Id, nil), h.Team)
		datastore.Put(t.context, datastore.NewKey(t.context, "Team", "", a.Id, nil), a.Team)
	}
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
	dsH := &ds.Team{}
	dsA := &ds.Team{}
	err = datastore.Get(
		t.context,
		datastore.NewKey(t.context, "Team", "", e.HomeId, nil),
		dsH)

	err = datastore.Get(
		t.context,
		datastore.NewKey(t.context, "Team", "", e.AwayId, nil),
		dsA)

	return &Team{Team: dsH, Id: e.HomeId}, &Team{Team: dsA, Id: e.AwayId}, err
}
