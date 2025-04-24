package main

import (
	"log"
	"os"

	App "id.benderaku.manufacture/app"
)

func main() {
	App.Initialize()
	if err := App.Run(":" + os.Getenv("SERVER_PORT")); err != nil {
		log.Fatal("Application failed: ", err)
	}
}
