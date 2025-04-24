package main

import (
	"os"

	"id.benderaku.manufacture/app"
)

func main() {
    a := app.App{}
    a.Initialize()
    if err := a.Run(":" + os.Getenv("SERVER_PORT")); err != nil {
		a.ERROR.Fatal("Application failed: ", err)
	}
}