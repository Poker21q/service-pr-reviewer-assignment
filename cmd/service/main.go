package main

import (
	"fmt"
	"os"

	"service-pr-reviewer-assignment/internal/app"
	"service-pr-reviewer-assignment/internal/app/config"
	"service-pr-reviewer-assignment/pkg/log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error(
			fmt.Sprintf("failed to load env %s", err.Error()),
		)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Error(
			fmt.Sprintf("cannot load config: %s", err.Error()),
		)
		os.Exit(1)
	}

	err = app.Run(cfg)
	if err != nil {
		log.Error(
			fmt.Sprintf("app exited with error: %s", err.Error()),
		)
		os.Exit(1)
	}
}
