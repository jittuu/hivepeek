package events

import (
	"encoding/csv"
	"io"
	"peek/ds"
	"strconv"
	"time"
)

func parseEvents(f io.Reader) ([]*ds.Event, error) {
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	events := make([]*ds.Event, len(lines)-1)

	for i, line := range lines[1:] {
		startTime, _ := time.Parse("02/01/06", line[1])
		hGoal, _ := strconv.ParseInt(line[4], 10, 32)
		aGoal, _ := strconv.ParseInt(line[5], 10, 32)
		mxH, _ := strconv.ParseFloat(line[54], 64)
		avH, _ := strconv.ParseFloat(line[55], 64)
		mxD, _ := strconv.ParseFloat(line[56], 64)
		avD, _ := strconv.ParseFloat(line[57], 64)
		mxA, _ := strconv.ParseFloat(line[58], 64)
		avA, _ := strconv.ParseFloat(line[59], 64)

		event := &ds.Event{
			StartTime: startTime,
			Home:      line[2],
			Away:      line[3],
			HGoal:     int(hGoal),
			AGoal:     int(aGoal),
			MaxOdds: ds.MatchOdds{
				Home: mxH,
				Draw: mxD,
				Away: mxA,
			},
			AvgOdds: ds.MatchOdds{
				Home: avH,
				Draw: avD,
				Away: avA,
			},
		}

		events[i] = event
	}

	return events, nil
}
