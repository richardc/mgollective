package mgollective

import (
	"github.com/simonz05/godis/redis"
	"launchpad.net/goyaml"
	"log"
	"regexp"
)

type RedisConnector struct {
	config *Config
	client *redis.Client
	subs   *redis.Sub
}

type RedisMessageWrapper struct {
	Headers map[string]string `yaml:":headers"`
	Body    string            `yaml:":body"`
}

func (r *RedisConnector) Connect() {
	log.Println("Connecting to redis")
}

func (r *RedisConnector) Subscribe() {
	var channels []string
	if !r.config.client {
		for _, collective := range r.config.collectives() {
			topic := collective + "::server::agents"
			channels = append(channels, topic)
		}
	} else {
		channels = append(channels, r.config.identity())
	}
	log.Println("Subscribing to ", channels)

	sub, err := r.client.Subscribe(channels...)
	if err != nil {
		log.Fatal(err)
	}
	r.subs = sub
}

func (r *RedisConnector) Publish(msg Message) {
	log.Printf("Publishing %+v", msg)
	body, err := goyaml.Marshal(&msg.body)
	if err != nil {
		log.Println("Failed to Marshal", err)
		return
	}

	var wrapper RedisMessageWrapper
	wrapper.Body = "---\n" + string(body) + "\n"
	headers := make(map[string]string, 0)
	if msg.reply_to != "" {
		headers["reply-to"] = msg.reply_to
	}
	wrapper.Headers = headers
	yaml_wrapper, err := goyaml.Marshal(&wrapper)
	if err != nil {
		log.Println("Failed to Marshal wrapper", err)
		return
	}
	log.Printf("Marshalled wrapper as %s", yaml_wrapper)

	r.client.Publish(msg.target, yaml_wrapper)
}

func (r *RedisConnector) Loop(parsed chan Message) {
	for msg := range r.subs.Messages {
		// ruby symbols/YAML encoding is special
		// Pretend like it was just a string with a colon
		silly_ruby, _ := regexp.Compile("!ruby/sym ")
		wire := silly_ruby.ReplaceAll(msg.Elem, []byte(":"))

		var wrapper RedisMessageWrapper
		if err := goyaml.Unmarshal(wire, &wrapper); err != nil {
			log.Println("YAML Unmarshal wrapper", err)
			continue
		}
		log.Printf("unpackged wrapper %+v", wrapper)

		var body MessageBody
		if err := goyaml.Unmarshal([]byte(wrapper.Body), &body); err != nil {
			log.Println("YAML Unmarshal body", err)
			continue
		}

		message := Message{
			topic:    msg.Channel,
			reply_to: wrapper.Headers["reply-to"],
			body:     body,
		}
		parsed <- message
	}
}

func makeRedisConnector(config *Config) Connector {
	log.Println("makeRedisConnector")
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
	registerConnector("redis", makeRedisConnector)
}
