package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseUrl string
	Port        uint
	JwtSecret   string
	Timeout     time.Duration
}

func New() *Config {
	DATABASE_URL := mustGetenv("DATABASE_URL")
	PORT := getenvOrDefaultParse("PORT", "8000", strconv.Atoi)
	JWT_SECRET := mustGetenv("JWT_SECRET")
	TIMEOUT := getenvOrDefaultParse("TIMEOUT", "5s", time.ParseDuration)

	return &Config{
		DATABASE_URL,
		uint(PORT),
		JWT_SECRET,
		TIMEOUT,
	}
}

func mustGetenv(k string) string {
	env, ok := os.LookupEnv(k)
	if !ok {
		slog.Warn("required env variable not provided",
			slog.String("variable", k),
		)
		panic(0)
	}
	return env
}

func getenvOrDefault(k string, v string) string {
	env, ok := os.LookupEnv(k)
	if !ok {
		slog.Warn("env variable not provided, using default value instead",
			slog.String("variable", k),
			slog.String("default", v),
		)
		return v
	}
	return env
}

func getenvOrDefaultParse[T any](k string, v string, parse func(string) (T, error)) T {
	envStr := getenvOrDefault(k, v)
	env, err := parse(envStr)
	if err != nil {
		slog.Error("failed to parse env variable",
			slog.String("variable", k),
			slog.String("value", envStr),
			slog.String("error", err.Error()),
		)
		panic(0)
	}
	return env
}
