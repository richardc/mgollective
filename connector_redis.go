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

func (r *RedisConnector) Loop(ch chan Message) {
	for {
		m := <-r.subs.Messages
		log.Println(m.Elem)

		// YAML Unmarshalling is wierd
		var message map[string]interface{}
		err := goyaml.Unmarshal([]byte(m.Elem), &message)
		if err != nil {
			log.Println("YAML Unmarshal", err)
		}
		log.Println("Headers ", message[":headers"])

		var body map[string]interface{}
		err = goyaml.Unmarshal([]byte(message[":body"].(string)), &body)

		ch <- Message{
			topic:   m.Channel,
			headers: message[":headers"],
			body:    body,
		}
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
