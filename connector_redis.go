package mgollective

import (
	"github.com/simonz05/godis/redis"
	"launchpad.net/goyaml"
	"log"
)

type RedisConnector struct {
	client *redis.Client
	subs   *redis.Sub
}

func (r *RedisConnector) Connect() {
	log.Println("Connecting to redis")
}

func (r *RedisConnector) Subscribe(config *Config) {
	var channels []string
	for _, collective := range config.collectives() {
		topic := collective + "::server::agents"
		channels = append(channels, topic)
	}
	log.Println("Subscribing to ", channels)

	sub, err := r.client.Subscribe(channels...)
	if err != nil {
		log.Fatal(err)
	}
	r.subs = sub
}

func (r *RedisConnector) Loop(parsed chan Message) {
	for msg := range r.subs.Messages {
		log.Println(msg.Elem)

		// YAML Unmarshalling is weird.  Extra weirded by the way these
		// messages are a document inside a document
		var wrapper struct {
			Headers map[string]string `yaml:":headers"`
			Body    string            `yaml:":body"`
		}
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
	return &RedisConnector{client: client}
}

func init() {
	registerConnector("redis", makeRedisConnector)
}
