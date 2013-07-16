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
	HRatingHome int
	HRatingAway int
	ARating     int
	ARatingHome int
	ARatingAway int
}

type Team struct {
	Name    string
	Season  string
	Overall int
	Home    int
	Away    int
}
