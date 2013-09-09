package mgollective

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var allowedConfig = map[string]bool{}

func DeclareConfig(name string) {
	allowedConfig[name] = true
}

func ParseConfig(file string) map[string]string {
	configValues := make(map[string]string)

	fh, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		key := fields[0]
		value := fields[2]
		if _, ok := allowedConfig[key]; ok {
			configValues[key] = value
		} else {
			// XXX Wuss out on rejecting for now
			//	log.Fataln("unexpected config key: '" + key + "'")
			log.Println("unexpected config key: '" + key + "'")
			configValues[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return configValues
}
