package mgollective

import (
    "log"
)

type RedisConnector struct {
}

func (r RedisConnector) Connect() int {
    log.Println("Connecting to redis")
	return 0
}

func redisConnect() Connector {
	return RedisConnector{}
}

func init() {
	registerConnector("redis", redisConnect)
}
