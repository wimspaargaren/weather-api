// Package main bootstraps the weather app.
package main

import (
	"log"

	"github.com/wimspaargaren/weather-api/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal("unexpected error while running the app", err)
	}
	log.Default().Println("app exited gracefully")
}
