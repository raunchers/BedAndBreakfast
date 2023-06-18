package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/raunchers/BedAndBreakfast/pkg/config"
	"github.com/raunchers/BedAndBreakfast/pkg/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	// create a file server, a place to get static files from
	fileServer := http.FileServer(http.Dir("./static/"))
	// Use the file server
	// Strip the URL
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
