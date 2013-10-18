package events

import (
	"strconv"
)

type CsvEvents []*Event

func (events CsvEvents) Csv() [][]string {
	header := []string{
		"League",
		"Season",
		"Start Time",
		"Home",
		"Away",
		"HG",
		"AG",
		"Rating",
		"AvgMatchOddsHome",
		"AvgMatchOddsDraw",
		"AvgMatchOddsAway",
		"MaxMatchOddsHome",
		"MaxMatchOddsDraw",
		"MaxMatchOddsAway",
		"AvgAHh",
		"AvgAHHome",
		"AvgAHAway",
		"MaxAHh",
		"MaxAHHome",
		"MaxAHAway",
		"AvgOUOver",
		"AvgOUUnder",
		"MaxOUOver",
		"MaxOUUnder",
	}

	csv := make([][]string, 0)

	csv = append(csv, header)

	for _, e := range events {
		row := []string{
			e.League,
			e.Season,
			e.StartTime.Format("02-01-2006"),
			e.Home,
			e.Away,
			strconv.Itoa(e.HGoal),
			strconv.Itoa(e.AGoal),
			strconv.FormatFloat(e.RatingDiff()/100, 'f', 2, 64),
			strconv.FormatFloat(e.AvgOdds.Home, 'f', 2, 64),
			strconv.FormatFloat(e.AvgOdds.Draw, 'f', 2, 64),
			strconv.FormatFloat(e.AvgOdds.Away, 'f', 2, 64),
			strconv.FormatFloat(e.MaxOdds.Home, 'f', 2, 64),
			strconv.FormatFloat(e.MaxOdds.Draw, 'f', 2, 64),
			strconv.FormatFloat(e.MaxOdds.Away, 'f', 2, 64),
			strconv.FormatFloat(e.AvgAHOdds.Handicap, 'f', 2, 64),
			strconv.FormatFloat(e.AvgAHOdds.Home, 'f', 2, 64),
			strconv.FormatFloat(e.AvgAHOdds.Away, 'f', 2, 64),
			strconv.FormatFloat(e.MaxAHOdds.Handicap, 'f', 2, 64),
			strconv.FormatFloat(e.MaxAHOdds.Home, 'f', 2, 64),
			strconv.FormatFloat(e.MaxAHOdds.Away, 'f', 2, 64),
			strconv.FormatFloat(e.AvgOUOdds.Over, 'f', 2, 64),
			strconv.FormatFloat(e.AvgOUOdds.Under, 'f', 2, 64),
			strconv.FormatFloat(e.MaxOUOdds.Over, 'f', 2, 64),
			strconv.FormatFloat(e.MaxOUOdds.Under, 'f', 2, 64),
		}

		csv = append(csv, row)
	}

	return csv
}
