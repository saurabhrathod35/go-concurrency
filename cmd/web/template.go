package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// set up middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		rw := w
		app.HomePage(rw, r)
	})

	return mux
}
