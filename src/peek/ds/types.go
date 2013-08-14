package ds

import (
	"time"
)

type Event struct {
	League      string
	Season      string
	StartTime   time.Time
	Home        string
	HomeId      int64
	Away        string
	AwayId      int64
	HGoal       int     `datastore:",noindex"`
	AGoal       int     `datastore:",noindex"`
	HRating     float64 `datastore:",noindex"`
	HNetRating  float64 `datastore:",noindex"`
	HFormRating float64 `datastore:",noindex"`
	ARating     float64 `datastore:",noindex"`
	ANetRating  float64 `datastore:",noindex"`
	AFormRating float64 `datastore:",noindex"`
	AvgOdds     MatchOdds
	MaxOdds     MatchOdds
}

type MatchOdds struct {
	Home float64 `datastore:",noindex"`
	Draw float64 `datastore:",noindex"`
	Away float64 `datastore:",noindex"`
}

type Team struct {
	Name                string
	Season              string
	OverallRating       float64
	HomeNetRating       float64
	AwayNetRating       float64
	LastFiveMatchRating []float64
}

func (t *Team) FormRating() float64 {
	total := 0.0
	for _, r := range t.LastFiveMatchRating {
		total += r
	}

	return total
}

func (t *Team) AddRating(rating float64) {
	ratings := t.LastFiveMatchRating
	if n := len(t.LastFiveMatchRating); n > 4 {
		ratings = t.LastFiveMatchRating[n-4:]
	}
	t.LastFiveMatchRating = append(ratings, rating)
}
