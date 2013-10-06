package events

import (
	"appengine"
	"appengine/datastore"
	"peek/ds"
	"peek/fx"
	"time"
)

const (
	epl        = 1204
	serieA     = 1269
	bundesliga = 1229
)

type fetchTask struct {
	context appengine.Context
	league  string
}

func (t *fetchTask) exec() error {
	s, err := fx.Fetch(t.context, t.league)
	if err != nil {
		return err
	}

	mappings, err := t.getTeamMappings()
	if err != nil {
		return err
	}

	season := getSeason(time.Now())
	teamIDs, err := t.getAllTeamIDs(season, mappings)
	if err != nil {
		return err
	}

	for _, l := range s.Leagues {
		if l.Id == epl || l.Id == serieA || l.Id == bundesliga {
			fixtures, fkeys, err := ds.GetFixturesByLeague(t.context, t.league)
			if err != nil {
				return err
			}

			visited := make(map[string]bool)
			for _, evt := range l.Events {
				// add missing mapping
				if !visited[evt.Home.Name] && mappings[evt.Home.Name] == "" {
					t.addTeamMapping(evt.Home.Name)
					visited[evt.Home.Name] = true
				}
				if !visited[evt.Away.Name] && mappings[evt.Away.Name] == "" {
					t.addTeamMapping(evt.Away.Name)
					visited[evt.Away.Name] = true
				}

				i, f := FixtureList(fixtures).Find(evt.Home.Name, evt.Away.Name, evt.StartTime())
				if f != nil {
					f.HomeId = int64(teamIDs[f.Home])
					f.AwayId = int64(teamIDs[f.Away])
					_, err = datastore.Put(t.context, fkeys[i], f)
				} else {
					f := &ds.Fixture{
						League:    t.league,
						Season:    season,
						StartTime: evt.StartTime(),
						Home:      evt.Home.Name,
						HomeId:    teamIDs[evt.Home.Name],
						Away:      evt.Away.Name,
						AwayId:    teamIDs[evt.Away.Name],
					}

					_, err = datastore.Put(
						t.context,
						datastore.NewIncompleteKey(t.context, "Fixture", nil),
						f)
				}
			}
		}
	}

	return nil
}

func (t *fetchTask) addTeamMapping(name string) error {
	m := &ds.TeamMapping{
		Name:       name,
		MasterName: name,
	}

	_, err := datastore.Put(
		t.context,
		datastore.NewIncompleteKey(t.context, "TeamMapping", nil),
		m)

	return err
}

func (t *fetchTask) getTeamMappings() (mappings map[string]string, err error) {
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

func (t *fetchTask) getAllTeamIDs(season string, teamMappings map[string]string) (teams map[string]int64, err error) {
	dst, keys, err := ds.GetAllTeams(t.context, season)
	if err != nil {
		return
	}

	teams = make(map[string]int64)
	for i, dt := range dst {
		if teamMappings[dt.Name] != "" {
			teams[teamMappings[dt.Name]] = keys[i].IntID()
		}
	}

	return
}

type FixtureList []*ds.Fixture

func (fxList FixtureList) Find(home, away string, startTime time.Time) (int, *ds.Fixture) {
	for i, f := range fxList {
		if f.Home == home && f.Away == away && f.StartTime.UTC() == startTime {
			return i, f
		}
	}

	return -1, nil
}
