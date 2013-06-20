package ds

import (
	"time"
)

type Event struct {
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

type Team struct {
	Name       string
	Rating     int
	RatingHome int
	RatingAway int
}
