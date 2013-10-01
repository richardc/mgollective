package mgollective

import (
	"fmt"
	"github.com/richardc/mgollective/mgollective"
	"time"
)

func pingAction(mgollective.RequestMessage) *mgollective.ResponseBody {
	return &mgollective.ResponseBody{
		"pong": fmt.Sprintf("%d", time.Now().Unix()),
	}
}

func init() {
	mgollective.RegisterAgent("rpcutil", mgollective.Agent{
		Actions: []mgollective.Action{
			{
				Name: "ping",
				Run:  pingAction,
			},
		},
	})
}
