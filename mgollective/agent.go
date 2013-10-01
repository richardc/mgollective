package mgollective

import (
	"github.com/golang/glog"
)

type Action struct {
	Name string
	Run  func(Mgollective, RequestMessage) *ResponseBody
}

type Agent struct {
	Actions []Action
}

var agentRegistry = map[string]Agent{}

func RegisterAgent(name string, agent Agent) {
	agentRegistry[name] = agent
}

func (a *Agent) Respond(app Mgollective, request RequestMessage) *ResponseMessage {
	for _, action := range a.Actions {
		if action.Name == request.Body.Action {
			body := action.Run(app, request)
			if body != nil {
				return &ResponseMessage{Body: *body}
			} else {
				glog.Infof("No response from %s", action.Name)
			}
		}
	}
	glog.Infof("No action called %s", request.Body.Action)

	return nil
}
