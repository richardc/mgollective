package mgollective

import (
	"github.com/msbranco/goconfig"
	"log"
)

type Config struct {
	config *goconfig.ConfigFile
}

func getconfig() *Config {
	conf, err := goconfig.ReadConfigFile("mgo.client.conf")
	if err != nil {
		log.Fatal("Couldn't open config ", err)
	}
	return &Config{config: conf}
}

func (c *Config) GetString(section, variable string) string {
	value, _ := c.config.GetString(section, variable)
	return value
}
