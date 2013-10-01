package mgo_agent_rpcutil

import (
	"fmt"
	"github.com/richardc/mgollective/mgollective"
	"time"
)

/// This is really a core agent, but putting it here proves the api is divorced

func pingAction(mgollective.Mgollective, mgollective.RequestMessage) *mgollective.ResponseBody {
	return &mgollective.ResponseBody{
		"pong": fmt.Sprintf("%d", time.Now().Unix()),
	}
}

func getFactAction(app mgollective.Mgollective, request mgollective.RequestMessage) *mgollective.ResponseBody {
	fact := request.Body.Params["fact"]
	return &mgollective.ResponseBody{
		"fact":  fact,
		"value": app.Factsource.GetFact(fact),
	}
}

func init() {
	mgollective.RegisterAgent("rpcutil", mgollective.Agent{
		Actions: []mgollective.Action{
			{Name: "ping", Run: pingAction},
			{Name: "get_fact", Run: getFactAction},
		},
	})
}
