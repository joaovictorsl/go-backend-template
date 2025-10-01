package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseUrl             string
	GoogleClientId          string
	GoogleClientSecret      string
	GoogleClientRedirectUrl string
	JwtSecret               string
	Port                    uint
	RequestTimeout          time.Duration
	ShutdownTimeout         time.Duration
	LogLevel                slog.Level
	AccessTokenTTL          time.Duration
	RefreshTokenTTL         time.Duration
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
		viper.GetString("google_client_id"),
		viper.GetString("google_client_secret"),
		viper.GetString("google_client_redirect_url"),
		viper.GetString("jwt_secret"),

		viper.GetUint("port"),
		viper.GetDuration("timeouts.request"),
		viper.GetDuration("timeouts.shutdown"),
		parseLogLevel(viper.GetString("log.level")),
		viper.GetDuration("token.access_ttl"),
		viper.GetDuration("token.refresh_ttl"),
	}
}

func setConfigDefaults() {
	viper.MustBindEnv("database_url")
	viper.MustBindEnv("google_client_id")
	viper.MustBindEnv("google_client_secret")
	viper.MustBindEnv("google_client_redirect_url")
	viper.MustBindEnv("jwt_secret")

	viper.SetDefault("api.port", 8000)
	viper.SetDefault("request.timeout", "10s")
	viper.SetDefault("shutdown.timeout", "1m")
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("token.access_ttl", "10m")
	viper.SetDefault("token.refresh_ttl", "168h")
}

func parseLogLevel(s string) slog.Level {
	var level slog.Level
	if err := level.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return level
}
