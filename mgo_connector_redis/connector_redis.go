package mgo_connector_redis

import (
	"github.com/richardc/mgollective/mgollective"
	"github.com/simonz05/godis/redis"
	"launchpad.net/goyaml"
	"regexp"
)

type RedisConnector struct {
	config *mgollective.Config
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
	if !r.config.IsClient() {
		for _, collective := range r.config.Collectives() {
			topic := collective + "::server::agents"
			channels = append(channels, topic)
		}
	} else {
		channels = append(channels, r.config.Identity())
	}
	mgollective.Logger().Debug("Subscribing to ", channels)

	sub, err := r.client.Subscribe(channels...)
	if err != nil {
		mgollective.Logger().Error(err)
		panic(err)
	}
	r.subs = sub
}

func (r *RedisConnector) Publish(msg mgollective.Message) {
	mgollective.Logger().Debugf("Publishing %+v", msg)
	body, err := goyaml.Marshal(&msg.Body)
	if err != nil {
		mgollective.Logger().Debugf("Failed to Marshal", err)
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
		mgollective.Logger().Debug("Failed to Marshal wrapper", err)
		return
	}
	mgollective.Logger().Tracef("Marshalled wrapper as %s", yaml_wrapper)

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
			mgollective.Logger().Debug("YAML Unmarshal wrapper", err)
			mgollective.Logger().Info("Recieved undecodable message, skipping it")
			continue
		}
		mgollective.Logger().Tracef("unpackged wrapper %+v", wrapper)

		var body mgollective.MessageBody
		if err := goyaml.Unmarshal([]byte(wrapper.Body), &body); err != nil {
			mgollective.Logger().Debug("YAML Unmarshal body", err)
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

func makeRedisConnector(config *mgollective.Config) mgollective.Connector {
	host := config.GetStringDefault("connector", "host", "127.0.0.1")
	port := config.GetStringDefault("connector", "port", "6379")
	db := config.GetIntDefault("connector", "database", 1)
	password := config.GetStringDefault("connector", "password", "")
	client := redis.New("tcp:"+host+":"+port, db, password)
	return &RedisConnector{
		config: config,
		client: client,
	}
}

func init() {
	mgollective.RegisterConnector("redis", makeRedisConnector)
}
