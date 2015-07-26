package backoffice

import (
	"net/http"
	"strconv"

	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"

	"github.com/gorilla/mux"
	"github.com/jittuu/hivepeek/internal"

	"golang.org/x/net/context"
)

var (
	delayFetchMatchesByLeagueFunc = delay.Func("fetch-matches-by-league", delayFetchMatchesByLeague)
)

func AllMatchesByLeague(c context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	l, err := strconv.Atoi(vars["league"])
	if err != nil {
		return err
	}
	s := vars["season"]
	db := &internal.DSContext{c}
	matches, err := db.GetAllMatchesByLeagueAndSeason(l, s)
	if err != nil {
		return err
	}

	return internal.Json(w, matches)
}

func FetchMatchesByLeague(c context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	l, err := strconv.Atoi(vars["league"])
	if err != nil {
		return err
	}
	s := vars["season"]
	return delayFetchMatchesByLeagueFunc.Call(c, l, s)
}

func delayFetchMatchesByLeague(c context.Context, league int, season string) error {
	client := Client(c)
	mths, err := client.GetFixturesByLeagueAndSeason(strconv.Itoa(league), season)
	log.Infof(c, "got %d matches from xmlsoccer for league: %d, season: %s", len(mths), league, season)
	if err != nil {
		return err
	}

	db := &internal.DSContext{c}
	existingMatches, err := db.GetAllMatchesByLeagueAndSeason(league, season)
	log.Infof(c, "have %d existing matches in datastore for league: %d, season: %s", len(existingMatches), league, season)
	if err != nil {
		return err
	}

	leagueName := GetLeagueNameByProviderID(db, league)
	var newMatches []*internal.Match
	for _, m := range mths {
		found := false
		for _, xm := range existingMatches {
			if xm.ProviderID == m.ID {
				found = true
				break
			}
		}

		if !found {
			newMatches = append(newMatches, &internal.Match{
				StartDate:        m.StartDate,
				Round:            m.Round,
				Status:           m.Time,
				HomeTeamID:       m.HomeTeamID,
				HomeTeamName:     m.HomeTeamName,
				HomeGoals:        m.HomeGoals,
				AwayTeamID:       m.AwayTeamID,
				AwayTeamName:     m.AwayTeamName,
				AwayGoals:        m.AwayGoals,
				ProviderID:       m.ID,
				ProviderName:     "xmlsoccer",
				LeagueProviderID: league,
				LeagueName:       leagueName,
				Season:           season,
			})
		}
	}

	log.Infof(c, "will add %d matches to datastore", len(newMatches))
	if len(newMatches) > 0 {
		err = db.PutMultiMatches(newMatches)
		if err != nil {
			return err
		}
	}

	return nil
}
