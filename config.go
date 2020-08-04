package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Dashboard server
type Dashboard struct {
	Name         string `yaml:"name"`
	URL          string `yaml:"url"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	IsThirdParty bool   `yaml:"is_third_party"`
	Seed         string `yaml:"seed"`
}

// Config to setup all confidential data
type Config struct {
	Dashboards []Dashboard `yaml:"dashboards"`
}

var cfgFile string = ".creds/config.yml"

func (cfg *Config) Read() {
	dat, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		log.Println(err)
		log.Fatalln("cannot read config file at " + cfgFile)
	}
	err = yaml.Unmarshal(dat, cfg)
	if err != nil {
		log.Println(err)
		log.Fatalln("cannot unmarshal config file at " + cfgFile)
	}
}
