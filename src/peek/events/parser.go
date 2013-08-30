package events

import (
	"encoding/csv"
	"io"
	"peek/ds"
	"strconv"
	"time"
)

func parseEvents(f io.Reader) ([]*ds.Event, error) {
	r := csv.NewReader(f)
	r.TrailingComma = true
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	events := make([]*ds.Event, len(lines)-1)
	h := getHeaderIndex(lines[0])

	for i, line := range lines[1:] {
		startTime, _ := time.Parse("02/01/06", line[h.StartTime])
		hGoal, _ := strconv.ParseInt(line[h.HGoal], 10, 32)
		aGoal, _ := strconv.ParseInt(line[h.AGoal], 10, 32)
		mxH, _ := strconv.ParseFloat(line[h.MxH], 64)
		avH, _ := strconv.ParseFloat(line[h.AvH], 64)
		mxD, _ := strconv.ParseFloat(line[h.MxD], 64)
		avD, _ := strconv.ParseFloat(line[h.AvD], 64)
		mxA, _ := strconv.ParseFloat(line[h.MxA], 64)
		avA, _ := strconv.ParseFloat(line[h.AvA], 64)

		event := &ds.Event{
			StartTime: startTime,
			Home:      line[h.Home],
			Away:      line[h.Away],
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

func getHeaderIndex(header []string) *headerIndex {
	h := &headerIndex{}
	for i, col := range header {
		switch col {
		case "Date":
			h.StartTime = i
		case "HomeTeam":
			h.Home = i
		case "AwayTeam":
			h.Away = i
		case "FTHG":
			h.HGoal = i
		case "FTAG":
			h.AGoal = i
		case "BbMxH":
			h.MxH = i
		case "BbMxD":
			h.MxD = i
		case "BbMxA":
			h.MxA = i
		case "BbAvH":
			h.AvH = i
		case "BbAvD":
			h.AvD = i
		case "BbAvA":
			h.AvA = i
		}
	}

	return h
}

type headerIndex struct {
	StartTime     int
	Home, Away    int
	HGoal, AGoal  int
	MxH, MxD, MxA int
	AvH, AvD, AvA int
}
