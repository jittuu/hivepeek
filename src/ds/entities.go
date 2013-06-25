package ds

import (
	"time"
)

type Event struct {
	Season      string
	StartTime   time.Time
	Home        string
	HomeId      int64
	Away        string
	AwayId      int64
	HGoal       int
	AGoal       int
	HRating     int
	HRatingHome int
	HRatingAway int
	ARating     int
	ARatingHome int
	ARatingAway int
}

type SeasonRating struct {
	Season  string
	Overall int
	Home    int
	Away    int
}

type Team struct {
	Name    string
	Ratings []SeasonRating
}
