package app

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	EXE "id.benderaku.manufacture/app/helpers"
	R "id.benderaku.manufacture/app/routes"
)

func Initialize() error {
	if err := EXE.InitDB(); err != nil {
		return err
	}

	if err := EXE.Logs(); err != nil {
		return err
	}

	R.RegisterRoutes()
	return nil
}

func Run(addr string) error {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3040", "http://localhost:3080", "https://manufaktur-dev.benderaku.id"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "token"}),
		handlers.AllowCredentials(),
	)

	log.Println("Server starting on", addr)
	return http.ListenAndServe(addr, cors(R.Router))
}
