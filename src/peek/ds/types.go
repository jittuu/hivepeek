package ds

import (
	"time"
)

type Fixture struct {
	League    string
	Season    string
	StartTime time.Time
	Home      string
	HomeId    int64
	Away      string
	AwayId    int64
}

type TeamMapping struct {
	Name       string
	MasterName string
}

type Event struct {
	League         string
	Season         string
	StartTime      time.Time
	Home           string
	HomeId         int64
	Away           string
	AwayId         int64
	HGoal          int     `datastore:",noindex"`
	AGoal          int     `datastore:",noindex"`
	HRating        float64 `datastore:",noindex"`
	HRatingLen     int     `datastore:",noindex"`
	HNetRating     float64 `datastore:",noindex"`
	HNetRatingLen  int     `datastore:",noindex"`
	HFormRating    float64 `datastore:",noindex"`
	HFormRatingLen int     `datastore:",noindex"`
	ARating        float64 `datastore:",noindex"`
	ARatingLen     int     `datastore:",noindex"`
	ANetRating     float64 `datastore:",noindex"`
	ANetRatingLen  int     `datastore:",noindex"`
	AFormRating    float64 `datastore:",noindex"`
	AFormRatingLen int
	AvgOdds        MatchOdds
	MaxOdds        MatchOdds
	AvgAHOdds      AHOdds
	MaxAHOdds      AHOdds
}

type MatchOdds struct {
	Home float64 `datastore:",noindex"`
	Draw float64 `datastore:",noindex"`
	Away float64 `datastore:",noindex"`
}

type AHOdds struct {
	Home     float64 `datastore:",noindex"`
	Away     float64 `datastore:",noindex"`
	Handicap float64 `datastore:",noindex"`
}

type Team struct {
	Name                string
	Season              string
	OverallRating       float64
	OverallRatingLen    int
	HomeNetRating       float64
	HomeNetRatingLen    int
	AwayNetRating       float64
	AwayNetRatingLen    int
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
