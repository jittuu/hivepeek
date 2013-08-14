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
	HGoal       int
	AGoal       int
	HRating     float64
	HNetRating  float64
	HFormRating float64
	ARating     float64
	ANetRating  float64
	AFormRating float64
	AvgOdds     MatchOdds
	MaxOdds     MatchOdds
}

type MatchOdds struct {
	Home float64
	Draw float64
	Away float64
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
