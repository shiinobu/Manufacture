package app

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"id.benderaku.manufacture/app/helpers"
	"id.benderaku.manufacture/app/routes"
)

type App struct {
	Router   *mux.Router
	INFO     *log.Logger
	ERROR    *log.Logger
}

func (a *App) Initialize() {
	if err := helpers.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	loggers, err := helpers.Logs()
	if err != nil {
		log.Fatal("Failed to set up logging:", err)
	}

	a.INFO = loggers.INFO
	a.ERROR = loggers.ERROR
	a.Router = mux.NewRouter()
	routes.RegisterRoutes(a.Router, a.INFO, a.ERROR)
}

func (a *App) Run(addr string) error {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3040", "http://localhost:3080"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "token"}),
		handlers.AllowCredentials(),
	)
	err := http.ListenAndServe(addr, cors(a.Router))
	if err != nil {
		a.ERROR.Fatal("Server failed to start:", err)
	}
	a.INFO.Println("Server starting on", addr)
	return nil
}
