package conf

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/olebedev/config"
)

var topics []string

func ReadConfig() []string {
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	pl, err := cfg.List("topics")
	if err != nil {
		panic(err)
	}

	for _, j := range pl {
		switch name := j.(type) {
		case string:
			val := fmt.Sprintf("%v", j)
			topics = append(topics, val)
		default:
			log.Fatalf("wrong type! %v %T\n", name, name)
		}
	}
	return topics
}
