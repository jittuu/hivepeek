package internal

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type DSContext struct {
	context.Context
}

func (c *DSContext) GetLeagueByProviderID(id int) (*League, error) {
	q := datastore.NewQuery(KindLeague).
		Filter("ProviderID =", id).
		Limit(1)
	var lgs []*League
	keys, err := q.GetAll(c, &lgs)
	if err != nil {
		return nil, err
	}

	if len(lgs) > 0 {
		lgs[0].ID = keys[0].IntID()
		return lgs[0], nil
	}
	return nil, nil
}

func (c *DSContext) GetAllLeagues() ([]*League, error) {
	lgs := make([]*League, 0)
	q := datastore.NewQuery(KindLeague)
	keys, err := q.GetAll(c, &lgs)
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
			keys[i] = datastore.NewIncompleteKey(c, KindLeague, nil)
		} else {
			keys[i] = datastore.NewKey(c, KindLeague, "", m.ID, nil)
		}
	}

	keys, err := datastore.PutMulti(c, keys, leagues)
	if err != nil {
		return err
	}
	for i := range leagues {
		leagues[i].ID = keys[i].IntID()
	}

	return nil
}

func (c *DSContext) GetAllTeamsByLeagueAndSeason(leagueID int, season string) ([]*Team, error) {
	teams := make([]*Team, 0)
	q := datastore.NewQuery(KindTeam).
		Filter("LeagueProviderID =", leagueID).
		Filter("Season =", season)
	keys, err := q.GetAll(c, &teams)
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
			keys[i] = datastore.NewIncompleteKey(c, KindTeam, nil)
		} else {
			keys[i] = datastore.NewKey(c, KindTeam, "", m.ID, nil)
		}
	}

	keys, err := datastore.PutMulti(c, keys, teams)
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

	matches := make([]*Match, 0)
	keys, err := q.GetAll(c, &matches)
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
			keys[i] = datastore.NewIncompleteKey(c, KindMatch, nil)
		} else {
			keys[i] = datastore.NewKey(c, KindMatch, "", m.ID, nil)
		}
	}

	keys, err := datastore.PutMulti(c, keys, matches)
	if err != nil {
		return err
	}
	for i := range matches {
		matches[i].ID = keys[i].IntID()
	}

	return nil
}
