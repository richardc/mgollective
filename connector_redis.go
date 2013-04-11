package mgollective

import (
	"github.com/simonz05/godis/redis"
	"launchpad.net/goyaml"
	"log"
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

func (r *RedisConnector) Publish(msg map[string]interface{}) {
	log.Printf("Publishing %+v", msg)
	target := msg["target"]
	reply_to := msg["reply-to"]
	delete(msg, "target")
	delete(msg, "reply-to")
	body, err := goyaml.Marshal(&msg)
	if err != nil {
		log.Println("Failed to Marshal", err)
		return
	}
	log.Printf("Marshalled body to %s", body)
	var wrapper RedisMessageWrapper
	wrapper.Body = "---\n" + string(body)
	headers := make(map[string]string, 0)
	if reply_to != nil {
		headers["reply-to"] = reply_to.(string)
	}
	wrapper.Headers = headers
	yaml_wrapper, err := goyaml.Marshal(&wrapper)
	if err != nil {
		log.Println("Failed to Marshal wrapper", err)
		return
	}
	log.Printf("Marshalled wrapper as %s", yaml_wrapper)

	r.client.Publish(target.(string), yaml_wrapper)
}

func (r *RedisConnector) Loop(parsed chan Message) {
	for msg := range r.subs.Messages {
		log.Println(msg.Elem)

		// YAML Unmarshalling is weird.  Extra weirded by the way these
		// messages are a document inside a document
		var wrapper RedisMessageWrapper
		err := goyaml.Unmarshal([]byte(msg.Elem), &wrapper)
		if err != nil {
			log.Println("YAML Unmarshal message", err)
		}

		// Unpack the :body key
		var message Message
		err = goyaml.Unmarshal([]byte(wrapper.Body), &message)
		if err != nil {
			log.Println("YAML Unmarshal body", err)
		}

		message.topic = msg.Channel
		message.reply_to = wrapper.Headers["reply-to"]
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
