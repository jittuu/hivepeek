package backoffice

import (
	"net/http"

	"github.com/jittuu/hivepeek/internal"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"

	"golang.org/x/net/context"
)

var (
	delayFetchLeaguesFunc = delay.Func("fetch-leagues", delayFetchCompetitions)
)

func AllLeagues(c context.Context, w http.ResponseWriter, r *http.Request) error {
	db := &internal.DSContext{c}
	lgs, err := db.GetAllLeagues()
	if err != nil {
		return err
	}
	return internal.Json(w, lgs)
}

func FetchLeagues(c context.Context, w http.ResponseWriter, r *http.Request) error {
	err := delayFetchLeaguesFunc.Call(c)
	return err
}

func delayFetchCompetitions(c context.Context) error {
	client := Client(c)
	xsLeagues, err := client.GetAllLeagues()
	if err != nil {
		return err
	}
	log.Infof(c, "got %d leagues from xmlsoccer", len(xsLeagues))

	db := &internal.DSContext{c}
	existingLeagues, err := db.GetAllLeagues()
	log.Infof(c, "have %d existing leagues in datastore", len(existingLeagues))
	if err != nil {
		return err
	}

	var newLeagues []*internal.League
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
		}
	}

	log.Infof(c, "will add %d leagues to datastore", len(newLeagues))
	if len(newLeagues) > 0 {
		err = db.PutMultiLeagues(newLeagues)
		if err != nil {
			return err
		}
	}

	return nil
}
