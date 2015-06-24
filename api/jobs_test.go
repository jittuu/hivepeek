package api

import (
	"testing"

	xs "github.com/jittuu/xmlsoccer"
	"github.com/mjibson/goon"

	"appengine/aetest"
	"appengine/datastore"
)

func TestPullLeagues(t *testing.T) {
	// prepare context
	c, err := aetest.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// mock the xmlsoccer response
	f := func() ([]*xs.League, error) {
		return []*xs.League{
			&xs.League{
				ID:   1,
				Name: "L1",
			},
			&xs.League{
				ID:   2,
				Name: "L2",
			},
		}, nil
	}

	// act
	err = pullLeaguesCore(c, f)
	if err != nil {
		t.Error(err)
	}

	// assert
	g := goon.FromContext(c)
	var leagues []*League
	_, err = g.GetAll(datastore.NewQuery("League"), &leagues)

	if err != nil {
		t.Error(err)
	}

	if len(leagues) != 2 {
		t.Errorf("expected %d leagues, got %d", 2, len(leagues))
	}
}

func TestPullLeaguesDelta(t *testing.T) {
	// prepare context
	c, err := aetest.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// prepare existing data
	g := goon.FromContext(c)
	existingLeague := &League{
		Name:         "L1",
		ProviderID:   1,
		ProviderName: "xmlsoccer",
	}
	_, err = g.Put(existingLeague)
	if err != nil {
		t.Error(err)
	}

	// mock the xmlsoccer response
	f := func() ([]*xs.League, error) {
		return []*xs.League{
			&xs.League{
				ID:   1,
				Name: "L1",
			},
			&xs.League{
				ID:   2,
				Name: "L2",
			},
		}, nil
	}

	// act
	err = pullLeaguesCore(c, f)
	if err != nil {
		t.Error(err)
	}

	// assert
	var leagues []*League
	_, err = g.GetAll(datastore.NewQuery("League"), &leagues)

	if err != nil {
		t.Error(err)
	}

	if len(leagues) != 2 {
		t.Errorf("expected %d leagues, got %d", 2, len(leagues))
	}
}
