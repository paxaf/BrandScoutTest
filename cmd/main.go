package main

import (
	"log"

	"github.com/paxaf/BrandScoutTest/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("failed creating app: %v", err)
	}
	if err = app.Run(); err != nil {
		log.Fatalf("error running app: %v", err)
	}
	if err = app.Close(); err != nil {
		log.Fatalf("error graceful shutdown: %v", err)
	}
}
