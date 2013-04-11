package mgollective

import (
	"github.com/simonz05/godis/redis"
	"log"
)

type RedisConnector struct {
	client *redis.Client
}

func (r *RedisConnector) Connect() {
	log.Println("Connecting to redis")
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
