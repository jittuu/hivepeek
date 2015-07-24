package backoffice

import (
	"net/http"
	"os"

	"github.com/jittuu/hivepeek/internal"
	"github.com/jittuu/xmlsoccer"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"golang.org/x/net/context"
)

var (
	delayFetchLeaguesFunc = delay.Func("fetch-leagues", delayFetchCompetitions)
)

func AllLeagues(c context.Context, w http.ResponseWriter, r *http.Request) error {
	var existingLeagues []*internal.League
	q := datastore.NewQuery(internal.KindLeague)
	keys, err := q.GetAll(c, &existingLeagues)
	if err != nil {
		return err
	}

	for i := range existingLeagues {
		existingLeagues[i].ID = keys[i].IntID()
	}
	return internal.Json(w, existingLeagues)
}

func FetchLeagues(c context.Context, w http.ResponseWriter, r *http.Request) error {
	err := delayFetchLeaguesFunc.Call(c)
	return err
}

func delayFetchCompetitions(c context.Context) error {
	client := &xmlsoccer.Client{
		BaseURL: xmlsoccer.DemoURL,
		APIKey:  os.Getenv("XMLSOCCER_API_KEY"),
		Client:  urlfetch.Client(c),
	}
	xsLeagues, err := client.GetAllLeagues()
	if err != nil {
		return err
	}
	log.Infof(c, "got %d leagues from xmlsoccer", len(xsLeagues))

	var existingLeagues []*internal.League
	q := datastore.NewQuery(internal.KindLeague)
	_, err = q.GetAll(c, &existingLeagues)
	log.Infof(c, "have %d existing leagues in datastore", len(existingLeagues))
	if err != nil {
		return err
	}

	var newLeagues []*internal.League
	var newLeagueKeys []*datastore.Key

	for _, l := range xsLeagues {
		found := false
		for _, xl := range existingLeagues {
			if xl.ProviderID == l.ID {
				found = true
				break
			}
		}
		if !found {
			newLeagues = append(newLeagues, &internal.League{
				Name:            l.Name,
				ProviderID:      l.ID,
				ProviderName:    "xmlsoccer",
				LatestMatchTime: l.LatestMatch,
			})
			newLeagueKeys = append(newLeagueKeys, datastore.NewIncompleteKey(c, internal.KindLeague, nil))
		}
	}

	log.Infof(c, "will add %d leagues to datastore", len(newLeagues))
	_, err = datastore.PutMulti(c, newLeagueKeys, newLeagues)
	if err != nil {
		return err
	}

	return nil
}
