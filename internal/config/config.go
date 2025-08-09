package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseUrl     string
	Port            uint
	JwtSecret       string
	RequestTimeout  time.Duration
	ShutdownTimeout time.Duration
	LogLevel        slog.Level
}

func New() *Config {
	DATABASE_URL := mustGetenv("DATABASE_URL")
	PORT := getenvOrDefaultParse("PORT", "8000", strconv.Atoi)
	JWT_SECRET := mustGetenv("JWT_SECRET")
	REQUEST_TIMEOUT := getenvOrDefaultParse("REQUEST_TIMEOUT", "10s", time.ParseDuration)
	SHUTDOWN_TIMEOUT := getenvOrDefaultParse("SHUTDOWN_TIMEOUT", "1m", time.ParseDuration)
	LOG_LEVEL := getenvOrDefaultParse("LOG_LEVEL", "debug", parseLogLevel)

	return &Config{
		DATABASE_URL,
		uint(PORT),
		JWT_SECRET,
		REQUEST_TIMEOUT,
		SHUTDOWN_TIMEOUT,
		LOG_LEVEL,
	}
}

func parseLogLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}

func mustGetenv(k string) string {
	env, ok := os.LookupEnv(k)
	if !ok {
		panic(fmt.Sprintf("required env variable not provided: %s", k))
	}
	return env
}

func getenvOrDefault(k string, v string) string {
	env, ok := os.LookupEnv(k)
	if !ok {
		fmt.Printf("%s not provided, using: %s\n", k, v)
		return v
	}
	return env
}

func getenvOrDefaultParse[T any](k string, v string, parse func(string) (T, error)) T {
	envStr := getenvOrDefault(k, v)
	env, err := parse(envStr)
	if err != nil {
		panic(fmt.Sprintf("failed to parse value for %s: %s\n%s\n", k, v, err))
	}
	return env
}
