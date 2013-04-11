package mgollective

import (
	"github.com/simonz05/godis/redis"
	"testing"
)

func TestLoop(t *testing.T) {
	in := make(chan *redis.Message)
	out := make(chan Message)
	connector := &RedisConnector{
		subs: &redis.Sub{Messages: in},
	}

	go connector.Loop(out)
	in <- &redis.Message{
		Channel: "mcollective::server::agents",
		Elem:    []byte(`---
:headers:
  reply-to: mcollective::reply::middleware.example.net::4004
:body: |
  ---
  :agent: discovery
  :filter:
    identity: []

    agent: []

    fact: []

    compound: []

    cf_class: []

  :senderid: middleware.example.net
  :ttl: 60
  :msgtime: 1365693540
  :collective: mcollective
  :requestid: 74ba247376f1518aa471ee61ff5f8245
  :callerid: user=vagrant
  :body: ping
`),
	}
	close(in)
	parsed := <-out
	t.Log("Parsed ", parsed)
}
