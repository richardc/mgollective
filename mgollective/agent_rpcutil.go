package mgollective

import (
	"fmt"
	"github.com/golang/glog"
	"time"
)

type RpcUtilAgent struct {
	app *Mgollective
}

func (agent *RpcUtilAgent) Respond(msg RequestMessage) *ResponseMessage {
	glog.Infof("RpcUtil agent handling %+v", msg)
	if msg.Body.Action == "ping" {
		return &ResponseMessage{
			Body: ResponseBody{
				"pong": fmt.Sprintf("%d", time.Now().Unix()),
			},
		}
	} else {
		return nil
	}
}

func makeRpcUtilAgent(m *Mgollective) Agent {
	return &RpcUtilAgent{app: m}
}

func init() {
	registerAgent("rpcutil", makeRpcUtilAgent)
}
