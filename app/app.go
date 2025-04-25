package app

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	EXE "id.benderaku.manufacture/app/helpers"
	R "id.benderaku.manufacture/app/routes"
)

func Initialize() {
	if err := EXE.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	if err := EXE.Logs(); err != nil {
		log.Fatal("Failed to set up logging:", err)
	}

	R.RegisterRoutes()
}

func Run(addr string) error {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3040", "http://localhost:3080", "https://manufaktur-dev.benderaku.id"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "token"}),
		handlers.AllowCredentials(),
	)
	err := http.ListenAndServe(addr, cors(R.Router))
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
	log.Println("Server starting on", addr)
	return nil
}
