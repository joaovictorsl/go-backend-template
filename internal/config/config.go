package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DATABASE_URL               string
	PORT                       uint
	GOOGLE_CLIENT_ID           string
	GOOGLE_CLIENT_SECRET       string
	GOOGLE_CLIENT_REDIRECT_URI string
	JWT_SECRET                 string
	JWT_ISS                    string
	ACCESS_TOKEN_EXP           time.Duration
	REFRESH_TOKEN_EXP          time.Duration
	TIMEOUT                    time.Duration
}

func must(name string, value string) {
	if value == "" {
		panic(fmt.Sprintf("env variable %s not provided", name))
	}
}

func NewConfig() *Config {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	must("DATABASE_URL", DATABASE_URL)

	PORT, err := strconv.ParseUint(os.Getenv("PORT"), 10, 0)
	if err != nil {
		log.Printf("failed to read env variable PORT, using default value instead: %s", err.Error())
		PORT = 8080
	}

	GOOGLE_CLIENT_ID := os.Getenv("GOOGLE_CLIENT_ID")
	must("GOOGLE_CLIENT_ID", GOOGLE_CLIENT_ID)
	GOOGLE_CLIENT_SECRET := os.Getenv("GOOGLE_CLIENT_SECRET")
	must("GOOGLE_CLIENT_SECRET", GOOGLE_CLIENT_SECRET)
	GOOGLE_CLIENT_REDIRECT_URI := os.Getenv("GOOGLE_CLIENT_REDIRECT_URI")
	must("GOOGLE_CLIENT_REDIRECT_URI", GOOGLE_CLIENT_REDIRECT_URI)

	JWT_SECRET := os.Getenv("JWT_SECRET")
	must("JWT_SECRET", JWT_SECRET)
	JWT_ISS := os.Getenv("JWT_ISS")
	must("JWT_ISS", JWT_ISS)

	var ACCESS_TOKEN_EXP time.Duration
	ACCESS_TOKEN_EXP_INT, err := strconv.ParseUint(os.Getenv("ACCESS_TOKEN_EXP"), 10, 0)
	if err != nil {
		log.Printf("failed to read env variable ACCESS_TOKEN_EXP, using default value instead: %s", err.Error())
		ACCESS_TOKEN_EXP = 15 * time.Minute
	} else {
		ACCESS_TOKEN_EXP = time.Duration(ACCESS_TOKEN_EXP_INT) * time.Second
	}

	var REFRESH_TOKEN_EXP time.Duration
	REFRESH_TOKEN_EXP_INT, err := strconv.ParseUint(os.Getenv("REFRESH_TOKEN_EXP"), 10, 0)
	if err != nil {
		log.Printf("failed to read env variable REFRESH_TOKEN_EXP, using default value instead: %s", err.Error())
		oneWeek := 7 * 24 * time.Hour
		REFRESH_TOKEN_EXP = oneWeek
	} else {
		REFRESH_TOKEN_EXP = time.Duration(REFRESH_TOKEN_EXP_INT) * time.Second
	}

	var TIMEOUT time.Duration
	TIMEOUT_INT, err := strconv.ParseUint(os.Getenv("TIMEOUT"), 10, 0)
	if err != nil {
		log.Printf("failed to read env variable TIMEOUT, using default value instead: %s", err.Error())
		TIMEOUT = 5 * time.Second
	} else {
		TIMEOUT = time.Duration(TIMEOUT_INT) * time.Second
	}

	return &Config{
		DATABASE_URL:               DATABASE_URL,
		PORT:                       uint(PORT),
		GOOGLE_CLIENT_ID:           GOOGLE_CLIENT_ID,
		GOOGLE_CLIENT_SECRET:       GOOGLE_CLIENT_SECRET,
		GOOGLE_CLIENT_REDIRECT_URI: GOOGLE_CLIENT_REDIRECT_URI,
		JWT_SECRET:                 JWT_SECRET,
		JWT_ISS:                    JWT_ISS,
		ACCESS_TOKEN_EXP:           ACCESS_TOKEN_EXP,
		REFRESH_TOKEN_EXP:          REFRESH_TOKEN_EXP,
		TIMEOUT:                    TIMEOUT,
	}
}
