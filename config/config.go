package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port       string
	IfPg       bool
	PgUser     string
	PgPassword string
	PgDatabase string
}

func ParseConfig(path string) *Config {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		log.Fatal("error parsing config:", err)
	}
	return &config
}
