package mgollective

import (
	"github.com/golang/glog"
)

type RpcUtilAgent struct {
	app *Mgollective
}

func (agent *RpcUtilAgent) Respond(msg RequestMessage) *ResponseMessage {
	glog.Infof("RpcUtil agent handling %+v", msg)
	if msg.Body.Action == "ping" {
		return &ResponseMessage{
			Body: ResponseBody{
				"message": "pong",
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
