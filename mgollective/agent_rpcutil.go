package mgollective

import (
	"github.com/golang/glog"
)

type RpcUtilAgent struct {
	app *Mgollective
}

func (a RpcUtilAgent) matches(msg RequestMessage) bool {
	return true
}

func (agent *RpcUtilAgent) Respond(msg RequestMessage) *ResponseMessage {
	glog.Infof("Discover agent handling %+v", msg)
	if !agent.matches(msg) {
		glog.Infof("Not for us")
		return nil
	}
	var body string
	if msg.Body.Action == "ping" {
		body = "pong"
	} else {
		body = "Unknown Request: " + msg.Body.Action
	}

	response := ResponseMessage{
		Body: map[string]string{"message": body},
	}

	return &response
}

func makeRpcUtilAgent(m *Mgollective) Agent {
	return &RpcUtilAgent{app: m}
}

func init() {
	registerAgent("rpcutil", makeRpcUtilAgent)
}
