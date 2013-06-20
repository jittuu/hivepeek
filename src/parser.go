package peek

import (
	"ds"
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

func ParseEvents(f io.Reader) ([]*ds.Event, error) {
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	events := make([]*ds.Event, len(lines)-1)

	for i, line := range lines[1:] {
		startTime, _ := time.Parse("02/01/06", line[1])
		hGoal, _ := strconv.ParseInt(line[4], 10, 32)
		aGoal, _ := strconv.ParseInt(line[5], 10, 32)
		event := &ds.Event{
			StartTime: startTime,
			Home:      line[2],
			Away:      line[3],
			HGoal:     int(hGoal),
			AGoal:     int(aGoal),
		}

		events[i] = event
	}

	return events, nil
}
