package api

import (
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"

	"appengine/delay"
)

type JobsApi struct{}

var (
	pullLeaguesFunc = delay.Func("pull-leagues", PullLeagues)
)

// PullLeaguesApi will create a delay task to pull leagues from provider (xmlsoccer)
// should handle at
// 	POST jobs/pullleagues
func (jobs *JobsApi) PullLeaguesApi(c endpoints.Context) {
	pullLeaguesFunc.Call(c)
}
