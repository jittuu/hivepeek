package ds

import (
	"time"
)

type Event struct {
	League     string
	Season     string
	StartTime  time.Time
	Home       string
	HomeId     int64
	Away       string
	AwayId     int64
	HGoal      int
	AGoal      int
	HRating    int
	HNetRating int
	ARating    int
	ANetRating int
	AvgOdds    MatchOdds
	MaxOdds    MatchOdds
}

type MatchOdds struct {
	Home float64
	Draw float64
	Away float64
}

type Team struct {
	Name          string
	Season        string
	OverallRating int
	HomeNetRating int
	AwayNetRating int
}
