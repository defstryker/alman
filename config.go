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

// Gmail API info
type Gmail struct {
	ID           string `yaml:"id"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RefreshToken string `yaml:"refresh_token"`
}

// Config to setup all confidential data
type Config struct {
	Dashboards []Dashboard `yaml:"dashboards"`
	Gmail      Gmail       `yaml:"gmail"`
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
