package xmlsoccer

import (
	"encoding/xml"
	"time"
)

type xmlroot struct {
	XMLName xml.Name `xml:"XMLSOCCER.COM"`
	Leagues []League `xml:"League"`
}

// League represent a soccer League
type League struct {
	// Id is unique identifier of a league
	ID int `xml:"Id"`

	// Name is the name of a league
	Name string

	// LatestMatch is the date of the last match for the league
	LatestMatch time.Time
}
