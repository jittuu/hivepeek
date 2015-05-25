package api

import (
	"github.com/jittuu/hivepeek/xmlsoccer"
	"github.com/mjibson/goon"

	"appengine"
	"appengine/datastore"
)

// PullLeagues is the handler of PullLeaguesApi's delay task
// It will get the data from xmlsoccer and insert into the datastore if not exists.
func PullLeagues(c appengine.Context) error {
	xs := xmlsoccer.DemoClient(config.APIKey)
	return pullLeaguesCore(c, xs.GetAllLeagues)
}

func pullLeaguesCore(c appengine.Context, f func() ([]*xmlsoccer.League, error)) error {
	xsLeagues, err := f()
	if err != nil {
		return err
	}
	c.Infof("got %d leagues from xmlsoccer", len(xsLeagues))

	g := goon.FromContext(c)
	var existingLeagues []*League
	_, err = g.GetAll(datastore.NewQuery("League"), &existingLeagues)
	c.Infof("have %d existing leagues in datastore", len(existingLeagues))
	if err != nil {
		return err
	}

	var newLeagues []*League

	for _, l := range xsLeagues {
		found := false
		for _, xl := range existingLeagues {
			if xl.ProviderID == l.ID {
				found = true
			}
		}
		if !found {
			newLeagues = append(newLeagues, &League{
				Name:         l.Name,
				ProviderID:   l.ID,
				ProviderName: "xmlsoccer",
			})
		}
	}

	c.Infof("will add %d leagues to datastore", len(newLeagues))
	_, err = g.PutMulti(newLeagues)
	if err != nil {
		return err
	}

	return nil
}
