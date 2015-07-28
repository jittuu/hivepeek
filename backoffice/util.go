package backoffice

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/net/context"

	"github.com/jittuu/hivepeek/internal"
	"github.com/jittuu/xmlsoccer"

	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
)

func GetLeagueNameByProviderID(c *internal.DSContext, leagueID int) string {
	key := fmt.Sprintf("leaguename:byproviderID:%d", leagueID)
	if item, err := memcache.Get(c, key); err == nil {
		return string(item.Value)
	} else if err == memcache.ErrCacheMiss {
		xleague, _ := c.GetLeagueByProviderID(leagueID)
		if xleague != nil {
			memcache.Set(c, &memcache.Item{
				Key:   key,
				Value: []byte(xleague.Name),
			})

			return xleague.Name
		}
	}
	return ""
}

func FillTeamsCache(c *internal.DSContext) error {
	teams, err := c.GetAllTeams()
	if err != nil {
		return err
	}
	items := make([]*memcache.Item, len(teams))
	for i := range teams {
		key := fmt.Sprintf("teamID:byproviderID:%d", teams[i].ProviderID)
		items[i] = &memcache.Item{
			Key:   key,
			Value: []byte(strconv.FormatInt(teams[i].ID, 10)),
		}
	}

	return memcache.AddMulti(c, items)
}

func GetTeamIDbyProviderID(c context.Context, providerID int) int64 {
	key := fmt.Sprintf("teamID:byproviderID:%d", providerID)
	id := int64(0)
	if item, err := memcache.Get(c, key); err == nil {
		id, _ = strconv.ParseInt(string(item.Value), 10, 0)
	}

	return id
}

func Client(c context.Context) *xmlsoccer.Client {
	return &xmlsoccer.Client{
		BaseURL: xmlsoccer.DemoURL,
		APIKey:  os.Getenv("XMLSOCCER_API_KEY"),
		Client:  urlfetch.Client(c),
	}
}
