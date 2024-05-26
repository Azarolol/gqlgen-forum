package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	IfPg       bool
	PgUser     string
	PgPassword string
	PgDatabase string
}

func ParseConfig(path string) *Config {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("error parsing config:", err)
	}
	ifPg := true
	if strings.ToLower(os.Getenv("IF_PG")) == "false" {
		ifPg = false
	}
	return &Config{
		Port:       os.Getenv("PORT"),
		IfPg:       ifPg,
		PgUser:     os.Getenv("PG_USER"),
		PgPassword: os.Getenv("PG_PASSWORD"),
		PgDatabase: os.Getenv("PG_DATABASE"),
	}
}
