package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Routes returns the router with all application routes
func (app *Application) Routes() http.Handler {
	r := chi.NewRouter()

	// Basic middleware
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// API routes
	r.Route("/api", func(r chi.Router) {
		// v1 API
		r.Route("/v1", func(r chi.Router) {
			// Public routes
			r.Post("/subscribe", app.SubscribeHandler)
			
			// Private admin routes (these should be protected in production)
			r.Route("/admin", func(r chi.Router) {
				r.Use(app.AdminAuthMiddleware)
				r.Get("/subscribers", app.GetAllSubscribersHandler)
			})
		})
	})

	return r
} 