package internal

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type DSContext struct {
	C context.Context
}

func (c *DSContext) GetAllLeagues() ([]*League, error) {
	var lgs []*League
	q := datastore.NewQuery(KindLeague)
	keys, err := q.GetAll(c.C, &lgs)
	if err != nil {
		return nil, err
	}

	for i := range lgs {
		lgs[i].ID = keys[i].IntID()
	}

	return lgs, nil
}

func (c *DSContext) PutMultiLeagues(leagues []*League) error {
	keys := make([]*datastore.Key, len(leagues))
	for i, m := range leagues {
		if m.ID == 0 {
			keys[i] = datastore.NewIncompleteKey(c.C, KindLeague, nil)
		} else {
			keys[i] = datastore.NewKey(c.C, KindLeague, "", m.ID, nil)
		}
	}

	keys, err := datastore.PutMulti(c.C, keys, leagues)
	if err != nil {
		return err
	}
	for i := range leagues {
		leagues[i].ID = keys[i].IntID()
	}

	return nil
}

func (c *DSContext) GetAllTeamsByLeagueAndSeason(leagueID int, season string) ([]*Team, error) {
	var teams []*Team
	q := datastore.NewQuery(KindTeam).
		Filter("LeagueProviderID =", leagueID).
		Filter("Season =", season)
	keys, err := q.GetAll(c.C, &teams)
	if err != nil {
		return nil, err
	}

	for i := range teams {
		teams[i].ID = keys[i].IntID()
	}

	return teams, nil
}

func (c *DSContext) PutMultiTeams(teams []*Team) error {
	keys := make([]*datastore.Key, len(teams))
	for i, m := range teams {
		if m.ID == 0 {
			keys[i] = datastore.NewIncompleteKey(c.C, KindTeam, nil)
		} else {
			keys[i] = datastore.NewKey(c.C, KindTeam, "", m.ID, nil)
		}
	}

	keys, err := datastore.PutMulti(c.C, keys, teams)
	if err != nil {
		return err
	}
	for i := range teams {
		teams[i].ID = keys[i].IntID()
	}

	return nil
}

func (c *DSContext) GetAllMatchesByLeagueAndSeason(leagueID int, season string) ([]*Match, error) {
	q := datastore.NewQuery(KindMatch).
		Filter("LeagueProviderID =", leagueID).
		Filter("Season =", season)

	var matches []*Match
	keys, err := q.GetAll(c.C, &matches)
	if err != nil {
		return nil, err
	}

	for i := range matches {
		matches[i].ID = keys[i].IntID()
	}

	return matches, nil
}

func (c *DSContext) PutMultiMatches(matches []*Match) error {
	keys := make([]*datastore.Key, len(matches))
	for i, m := range matches {
		if m.ID == 0 {
			keys[i] = datastore.NewIncompleteKey(c.C, KindMatch, nil)
		} else {
			keys[i] = datastore.NewKey(c.C, KindMatch, "", m.ID, nil)
		}
	}

	keys, err := datastore.PutMulti(c.C, keys, matches)
	if err != nil {
		return err
	}
	for i := range matches {
		matches[i].ID = keys[i].IntID()
	}

	return nil
}
