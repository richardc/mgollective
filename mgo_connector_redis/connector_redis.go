package mgo_connector_redis

import (
	"github.com/richardc/mgollective/mgollective"
	"github.com/simonz05/godis/redis"
	"launchpad.net/goyaml"
	"regexp"
	"strconv"
)

type RedisConnector struct {
	app    *mgollective.Mgollective
	client *redis.Client
	subs   *redis.Sub
}

type RedisMessageWrapper struct {
	Headers map[string]string `yaml:":headers"`
	Body    string            `yaml:":body"`
}

func (r *RedisConnector) Connect() {
	// for this connector it's a noop
}

func (r *RedisConnector) Subscribe() {
	var channels []string
	if !r.app.IsClient() {
		for _, collective := range r.app.Collectives() {
			topic := collective + "::server::agents"
			channels = append(channels, topic)
		}
	} else {
		channels = append(channels, r.app.Identity())
	}
	r.app.Debug("Subscribing to ", channels)

	sub, err := r.client.Subscribe(channels...)
	if err != nil {
		r.app.Error(err)
		panic(err)
	}
	r.subs = sub
}

func (r *RedisConnector) Publish(msg mgollective.Message) {
	r.app.Debugf("Publishing %+v", msg)
	body, err := goyaml.Marshal(&msg.Body)
	if err != nil {
		r.app.Debugf("Failed to Marshal", err)
		return
	}

	var wrapper RedisMessageWrapper
	wrapper.Body = "---\n" + string(body) + "\n"
	headers := make(map[string]string, 0)
	if msg.Reply_to != "" {
		headers["reply-to"] = msg.Reply_to
	}
	wrapper.Headers = headers
	yaml_wrapper, err := goyaml.Marshal(&wrapper)
	if err != nil {
		r.app.Debug("Failed to Marshal wrapper", err)
		return
	}
	r.app.Debugf("Marshalled wrapper as %s", yaml_wrapper)

	r.client.Publish(msg.Target, yaml_wrapper)
}

func (r *RedisConnector) Loop(parsed chan mgollective.Message) {
	for msg := range r.subs.Messages {
		// ruby symbols/YAML encoding is special
		// Pretend like it was just a string with a colon
		silly_ruby, _ := regexp.Compile("!ruby/sym ")
		wire := silly_ruby.ReplaceAll(msg.Elem, []byte(":"))

		var wrapper RedisMessageWrapper
		if err := goyaml.Unmarshal(wire, &wrapper); err != nil {
			r.app.Debug("YAML Unmarshal wrapper", err)
			r.app.Info("Recieved undecodable message, skipping it")
			continue
		}
		r.app.Tracef("unpackged wrapper %+v", wrapper)

		var body mgollective.MessageBody
		if err := goyaml.Unmarshal([]byte(wrapper.Body), &body); err != nil {
			r.app.Debug("YAML Unmarshal body", err)
			continue
		}

		message := mgollective.Message{
			Topic:    msg.Channel,
			Reply_to: wrapper.Headers["reply-to"],
			Body:     body,
		}
		parsed <- message
	}
}

func makeRedisConnector(app *mgollective.Mgollective) mgollective.Connector {
	host := app.GetConfig("plugin.redis.host", "127.0.0.1")
	port := app.GetConfig("plugin.redis.port", "6379")
	db, _ := strconv.Atoi(app.GetConfig("plugin.redis.database", "1"))
	password := app.GetConfig("plugin.redis.password", "")
	client := redis.New("tcp:"+host+":"+port, db, password)
	return &RedisConnector{
		app:    app,
		client: client,
	}
}

func init() {
	mgollective.RegisterConnector("redis", makeRedisConnector)
}
