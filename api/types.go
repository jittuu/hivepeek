package api

import (
	"time"
)

type League struct {
	Id              int64 `datastore:"-" goon:"id"`
	Name            string
	LatestMatchTime time.Time `datastore:",noindex"`
	ProviderID      int
	ProviderName    string
}

type Event struct {
	Id           int64 `datastore:"-" goon:"id"`
	StartDate    time.Time
	Round        int
	Status       string
	HomeTeamName string
	HomeTeamID   int
	HomeGoals    int `datastore:",noindex"`
	AwayTeamName string
	AwayTeamID   int
	AwayGoals    int `datastore:",noindex"`
	ProviderID   int
	ProviderName string
}

type Team struct {
	Id           int64 `datastore:"-" goon:"id"`
	Name         string
	Country      string
	ProviderID   int
	ProviderName string
}
