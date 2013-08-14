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
	HRating     int
	HNetRating  int
	HFormRating int
	ARating     int
	ANetRating  int
	AFormRating int
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
	OverallRating       int
	HomeNetRating       int
	AwayNetRating       int
	LastFiveMatchRating []int
}

func (t *Team) FormRating() int {
	total := 0
	for _, r := range t.LastFiveMatchRating {
		total += r
	}

	return total
}

func (t *Team) AddRating(rating int) {
	ratings := t.LastFiveMatchRating
	if n := len(t.LastFiveMatchRating); n > 4 {
		ratings = t.LastFiveMatchRating[n-4:]
	}
	t.LastFiveMatchRating = append(ratings, rating)
}
