package mgo_connector_redis

import (
	"github.com/richardc/mgollective/mgollective"
	"github.com/simonz05/godis/redis"
	. "launchpad.net/gocheck"
	"testing"
)

// Magic gocheck boilerplate http://labix.org/gocheck
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

// Send a YAML document, test that it came back in the structures we expected
func (s *MySuite) TestLoop(c *C) {
	in := make(chan *redis.Message)
	out := make(chan mgollective.Message)
	connector := &RedisConnector{
		subs: &redis.Sub{Messages: in},
	}

	go connector.Loop(out)
	in <- &redis.Message{
		Channel: "mcollective::server::agents",
		Elem: []byte(`---
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
	c.Logf("Parsed %#v", parsed)

	c.Check(parsed.Topic, Equals, "mcollective::server::agents")
	c.Check(parsed.Reply_to, Equals, "mcollective::reply::middleware.example.net::4004")
	c.Check(parsed.Body.Agent, Equals, "discovery")
	c.Check(parsed.Body.Senderid, Equals, "middleware.example.net")
	c.Check(parsed.Body.Collective, Equals, "mcollective")
	c.Check(parsed.Body.Callerid, Equals, "user=vagrant")
}
