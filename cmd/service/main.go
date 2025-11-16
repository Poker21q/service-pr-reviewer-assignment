package main

import (
	"context"
	"os"

	"service-pr-reviewer-assignment/internal/app"
	"service-pr-reviewer-assignment/internal/app/config"
	"service-pr-reviewer-assignment/pkg/log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Errorf("failed to load env %v", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Errorf("cannot load config: %v", err)
		os.Exit(1)
	}

	err = app.Run(context.Background(), cfg)
	if err != nil {
		log.Errorf("app exited with error: %v", err)
		os.Exit(1)
	}
}
