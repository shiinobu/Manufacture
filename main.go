package main

import (
	"log"
	"os"

	App "id.benderaku.manufacture/app"
	"id.benderaku.manufacture/app/config"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	if err := App.Initialize(); err != nil {
		log.Fatal("Failed to initialize application: ", err)
	}

	if err := App.Run(":" + os.Getenv("SERVER_PORT")); err != nil {
		log.Fatal("Application failed: ", err)
	}
}
