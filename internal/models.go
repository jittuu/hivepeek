package internal

import (
	"time"
)

const (
	KindLeague = "League"
	KindMatch  = "Match"
	KindTeam   = "Team"
)

type League struct {
	ID              int64 `datastore:"-"`
	Name            string
	LatestMatchTime time.Time `datastore:",noindex"`
	ProviderID      int
	ProviderName    string
}

type Match struct {
	ID               int64 `datastore:"-"`
	StartDate        time.Time
	Round            int
	Status           string
	HomeTeamName     string
	HomeTeamID       int
	HomeGoals        int `datastore:",noindex"`
	AwayTeamName     string
	AwayTeamID       int
	AwayGoals        int `datastore:",noindex"`
	ProviderID       int
	ProviderName     string
	LeagueProviderID int
	LeagueName       string
	Season           string
}

type Team struct {
	ID               int64 `datastore:"-"`
	Name             string
	WikiLink         string
	Country          string
	ProviderID       int
	ProviderName     string
	LeagueProviderID int
	LeagueName       string
	Season           string
}
