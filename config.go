package mgollective

import (
    "log"
    "github.com/msbranco/goconfig"
)

func getconfig() *goconfig.ConfigFile {
    conf, err := goconfig.ReadConfigFile("mgo.client.conf")
    if err != nil {
        log.Fatal("Couldn't open config ", err)
    }
    return conf
}
