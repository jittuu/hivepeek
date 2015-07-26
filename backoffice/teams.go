package backoffice

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jittuu/hivepeek/internal"

	"golang.org/x/net/context"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
)

var (
	delayFetchTeamsFunc = delay.Func("fetch-teams-by-league", delayFetchTeams)
)

func AllTeamsByLeague(c context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	l, err := strconv.Atoi(vars["league"])
	if err != nil {
		return err
	}
	s := vars["season"]
	db := &internal.DSContext{c}
	teams, err := db.GetAllTeamsByLeagueAndSeason(l, s)
	if err != nil {
		return err
	}

	return internal.Json(w, teams)
}

func FetchTeamsByLeague(c context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	l, err := strconv.Atoi(vars["league"])
	if err != nil {
		return err
	}
	s := vars["season"]
	return delayFetchTeamsFunc.Call(c, l, s)
}

func delayFetchTeams(c context.Context, league int, season string) error {
	client := Client(c)
	xsTeams, err := client.GetAllTeamsByLeagueAndSeason(strconv.Itoa(league), season)
	if err != nil {
		return err
	}
	log.Infof(c, "got %d teams from xmlsoccer for league: %d, season: %s", len(xsTeams), league, season)

	db := &internal.DSContext{c}
	existingTeams, err := db.GetAllTeamsByLeagueAndSeason(league, season)
	log.Infof(c, "have %d existing teams for league: %d, season: %s in datastore", len(existingTeams), league, season)
	if err != nil {
		return err
	}

	xleague, err := db.GetLeagueByProviderID(league)
	if err != nil {
		return err
	}
	leagueName := ""
	if xleague != nil {
		leagueName = xleague.Name
	}

	var newTeams []*internal.Team
	for _, t := range xsTeams {
		found := false
		for _, xt := range existingTeams {
			if xt.ProviderID == t.ID {
				found = true
				break
			}
		}
		if !found {
			newTeams = append(newTeams, &internal.Team{
				Name:             t.Name,
				Country:          t.Country,
				WikiLink:         t.WikiLink,
				ProviderID:       t.ID,
				ProviderName:     "xmlsoccer",
				LeagueProviderID: league,
				LeagueName:       leagueName,
				Season:           season,
			})
		}
	}

	log.Infof(c, "will add %d teams to datastore", len(newTeams))
	if len(newTeams) > 0 {
		err = db.PutMultiTeams(newTeams)
		if err != nil {
			return err
		}
	}

	return nil
}
