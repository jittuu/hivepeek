package ds

import (
	"time"
)

type EventGoals struct {
	Id             int64 `datastore:"-"`
	League, Season string
	StartTime      time.Time
	Home, Away     string
	HomeGoals      []int `datastore:",noindex"`
	AwayGoals      []int `datastore:",noindex"`
	EventId        int64
}

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
	HGoal          int         `datastore:",noindex"`
	HGoals         []int       `datastore:",noindex"`
	AGoals         []int       `datastore:",noindex"`
	AGoal          int         `datastore:",noindex"`
	HRating        float64     `datastore:",noindex"`
	HRatingLen     int         `datastore:",noindex"`
	HNetRating     float64     `datastore:",noindex"`
	HNetRatingLen  int         `datastore:",noindex"`
	HFormRating    float64     `datastore:",noindex"`
	HFormRatingLen int         `datastore:",noindex"`
	ARating        float64     `datastore:",noindex"`
	ARatingLen     int         `datastore:",noindex"`
	ANetRating     float64     `datastore:",noindex"`
	ANetRatingLen  int         `datastore:",noindex"`
	AFormRating    float64     `datastore:",noindex"`
	AFormRatingLen int         `datastore:",noindex"`
	HGoalsFor      EventSector `datastore:",noindex"`
	HGoalsAgainst  EventSector `datastore:",noindex"`
	AGoalsFor      EventSector `datastore:",noindex"`
	AGoalsAgainst  EventSector `datastore:",noindex"`
	AvgOdds        MatchOdds
	MaxOdds        MatchOdds
	AvgAHOdds      AHOdds
	MaxAHOdds      AHOdds
	AvgOUOdds      OUOdds
	MaxOUOdds      OUOdds
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

type OUOdds struct {
	Over     float64 `datastore:",noindex"`
	Under    float64 `datastore:",noindex"`
	Handicap float64 `datastore:",noindex"`
}

type EventSector struct {
	M0_15  int `datastore:",noindex"`
	M16_30 int `datastore:",noindex"`
	M31_45 int `datastore:",noindex"`
	M46_60 int `datastore:",noindex"`
	M61_75 int `datastore:",noindex"`
	M76_90 int `datastore:",noindex"`
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
	GoalsFor            EventSector `datastore:",noindex"`
	GoalsAgainst        EventSector `datastore:",noindex"`
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

func (t *Team) CalcEventGoalsSector(f, a []int) {
	tf := &t.GoalsFor
	for _, g := range f {
		switch {
		case g > 75:
			tf.M76_90 += 1
		case g > 60:
			tf.M61_75 += 1
		case g > 45:
			tf.M46_60 += 1
		case g > 30:
			tf.M31_45 += 1
		case g > 15:
			tf.M16_30 += 1
		default:
			tf.M0_15 += 1
		}
	}
	ta := &t.GoalsAgainst
	for _, g := range a {
		switch {
		case g > 75:
			ta.M76_90 += 1
		case g > 60:
			ta.M61_75 += 1
		case g > 45:
			ta.M46_60 += 1
		case g > 30:
			ta.M31_45 += 1
		case g > 15:
			ta.M16_30 += 1
		default:
			ta.M0_15 += 1
		}
	}
}
