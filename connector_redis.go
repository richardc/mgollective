package mgollective

import (
	"log"
)

type RedisConnector struct {
}

func (r *RedisConnector) Connect() {
	log.Println("Connecting to redis")
}

func makeRedisConnector(config *Config) Connector {
	log.Println("makeRedisConnector")
	host := "192.168.1.20"
	port := "6379"
	db := 1
	password := ""
	client := redis.New("tcp:"+host+":"+port, db, password)
	return RedisConnector{client: client}
}

func init() {
	registerConnector("redis", makeRedisConnector)
}
