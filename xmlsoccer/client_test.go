package xmlsoccer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAllLeagues(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := `
<?xml version="1.0" encoding="utf-8"?>
<XMLSOCCER.COM>
    <League>
        <Id>1</Id>
        <Name>English Premier League</Name>
        <Country>England</Country>
        <Historical_Data>Yes</Historical_Data>
        <Fixtures>Yes</Fixtures>
        <Livescore>Yes</Livescore>
        <NumberOfMatches>2557</NumberOfMatches>
        <LatestMatch>2013-03-02T16:00:00+01:00</LatestMatch>
    </League>
    <League>
        <Id>3</Id>
        <Name>Scottish Premier League</Name>
        <Country>Scotland</Country>
        <Historical_Data>Yes</Historical_Data>
        <Fixtures>Yes</Fixtures>
        <Livescore>Yes</Livescore>
        <NumberOfMatches>1314</NumberOfMatches>
        <LatestMatch>2013-03-02T16:00:00+01:00</LatestMatch>
    </League>
    <League>
        <Id>4</Id>
        <Name>Bundesliga</Name>
        <Country>Germany</Country>
        <Historical_Data>Yes</Historical_Data>
        <Fixtures>Yes</Fixtures>
        <Livescore>Yes</Livescore>
        <NumberOfMatches>1743</NumberOfMatches>
        <LatestMatch>2013-03-02T15:30:00+01:00</LatestMatch>
    </League>
    <AccountInformation>Data requested at 02-03-2013 21:02:09 from XX.XX.XX.XX, Username: Espectro. Your current supscription runs out on XX-XX-XXXX 11:01:25.</AccountInformation>
</XMLSOCCER.COM>
    `
		fmt.Fprintln(w, data)
	}))
	defer ts.Close()

	// act
	c := DemoClient("dummy-key")
	c.testURL = ts.URL
	leagues, err := c.GetAllLeagues()

	// assert
	if err != nil {
		t.Error(err)
	}
	if len(leagues) != 3 {
		t.Errorf("expected %d leagues but return %d", 3, len(leagues))
		t.Log(leagues)
	}

	epl := leagues[0]
	if epl.ID != 1 {
		t.Errorf("expected ID: %d, got %d", 1, epl.ID)
	}
	if epl.Name != "English Premier League" {
		t.Errorf("expected name: %q, got %q", "English Premier League", epl.Name)
	}
	expectedTime, _ := time.Parse("2006-01-02T15:04:00+00:00", "2013-03-02T16:00:00+01:00")
	if epl.LatestMatch == expectedTime {
		t.Errorf("expected time: %v, got %v", expectedTime, epl.LatestMatch)
	}
}
