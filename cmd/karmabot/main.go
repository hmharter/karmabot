package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Hireology/karmabot/internal/karmabot"
	"go.uber.org/zap"
)

func main() {
	// Check Environment
	checkEnv()

	// Setup logger config
	var cfg zap.Config
	if os.Getenv("ENVIRONMENT") == karmabot.Production {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// Build sugared logger
	zaplog, err := cfg.Build()
	if err != nil {
		log.Fatal(fmt.Errorf("can't initialize zap logger: %w", err))
	}
	defer zaplog.Sync()
	logger := zaplog.Sugar()

	// Start Karmabot
	logger.Info("starting karmabot")

	karmabot.Karma, err = karmabot.LoadKarma()
	if err != nil {
		logger.Error(fmt.Errorf("error loading karma: %w", err))
	}

	c := make(chan int, 2)
	for {
		c <- 1
		logger.Debug("attempting new connection")
		go karmabot.OpenConnection(c, logger)
	}
}

func checkEnv() {
	ok := true

	envVars := []string{
		"ENVIRONMENT",
		"SLACK_BOT_TOKEN",
		"SLACK_APP_TOKEN",
		"SLACK_BOT_USER",
		"KARMAFILE",
		"SLACK_CONN_URL",
		"SLACK_POST_MESSAGE_URL",
		"SLACK_POST_EPHEMERAL_MESSAGE_URL",
	}

	for _, key := range envVars {
		var val string
		val, ok = os.LookupEnv(key)
		if !ok || len(val) == 0 {
			ok = false
		}
	}

	if !ok {
		log.Fatalf("Missing required environment variable(s). Check README for info.")
	}
}
