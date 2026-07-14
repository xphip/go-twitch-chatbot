package main

import (
	"log"
	"os"
)

type config struct {
	Nick    string
	Channel string
	Token   string
}

func loadConfig() config {
	return config{
		Nick:    mustEnv("TWITCH_NICK"),
		Channel: mustEnv("TWITCH_CHANNEL"),
		Token:   mustEnv("TWITCH_TOKEN"),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is required", key)
	}
	return v
}
