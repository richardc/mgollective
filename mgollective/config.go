package mgollective

import (
	"github.com/msbranco/goconfig"
	"os"
)

type Config struct {
	config *goconfig.ConfigFile
	client bool
}

func getconfig(file string, client bool) *Config {
	conf, err := goconfig.ReadConfigFile(file)
	if err != nil {
		logger.Criticalf("Couldn't open config ", err)
		os.Exit(1)
	}
	return &Config{config: conf, client: client}
}

func (c *Config) Collectives() []string {
	return []string{"mcollective"}
}

func (c *Config) Classes() []string {
	return []string{"mgollective"}
}

func (c *Config) Identity() string {
	return "mcollective::agent::pies"
}

func (c *Config) Callerid() string {
	return "user=meat"
}

func (c *Config) Senderid() string {
	return "meat.example.com"
}

func (c *Config) IsClient() bool {
	return c.client
}

func (c *Config) GetStringDefault(section, variable, def string) string {
	value, err := c.config.GetString(section, variable)
	if err != nil {
		return def
	}
	return value
}

func (c *Config) GetIntDefault(section, variable string, def int) int {
	value, err := c.config.GetInt64(section, variable)
	if err != nil {
		return def
	}
	return int(value)
}