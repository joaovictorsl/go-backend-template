package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseUrl        string
	Port               uint
	RequestTimeout     time.Duration
	ShutdownTimeout    time.Duration
	LogLevel           slog.Level
	GoogleClientId     string
	GoogleClientSecret string
}

func init() {
	viper.SetConfigFile("config.yml")
	viper.AddConfigPath(".")
}

func New() *Config {
	setConfigDefaults()
	err := viper.ReadInConfig()
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	return &Config{
		viper.GetString("database_url"),
		viper.GetUint("api.port"),
		viper.GetDuration("request.timeout"),
		viper.GetDuration("shutdown.timeout"),
		parseLogLevel(viper.GetString("log.level")),
		viper.GetString("google_client_id"),
		viper.GetString("google_client_secret"),
	}
}

func setConfigDefaults() {
	viper.SetDefault("api.port", 8000)
	viper.SetDefault("request.timeout", "10s")
	viper.SetDefault("shutdown.timeout", "1m")
	viper.SetDefault("log.level", "debug")
	viper.AutomaticEnv()
}

func parseLogLevel(s string) slog.Level {
	var level slog.Level
	if err := level.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return level
}
