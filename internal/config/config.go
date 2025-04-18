package config

import (
	"fmt"
	"os"
)

type Config struct {
	DATABASE_URL string
}

func must(name string, value string) {
	if value == "" {
		panic(fmt.Sprintf("env variable %s not provided", name))
	}
}

func NewConfig() *Config {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	must("DATABASE_URL", DATABASE_URL)

	return &Config{
		DATABASE_URL: DATABASE_URL,
	}
}
