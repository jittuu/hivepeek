package backoffice

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	"github.com/jittuu/hivepeek/internal"
	"github.com/jittuu/xmlsoccer"

	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
)

func GetLeagueNameByProviderID(c *internal.DSContext, leagueID int) string {
	key := fmt.Sprintf("leaguename:byproviderid:%d", leagueID)
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

func Client(c context.Context) *xmlsoccer.Client {
	return &xmlsoccer.Client{
		BaseURL: xmlsoccer.DemoURL,
		APIKey:  os.Getenv("XMLSOCCER_API_KEY"),
		Client:  urlfetch.Client(c),
	}
}
