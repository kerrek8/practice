package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(s.AuthMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", s.healthHandler)

	fs := http.StripPrefix("/", http.FileServer(http.Dir("front/")))
	r.Handle("/*", fs)

	r.Post("/api/register", s.RegisterHandler)
	r.Post("/api/login", s.LoginHandler)
	r.Post("/api/logout", s.LogoutHandler)
	r.Get("/api/me", s.MeHandler)
	r.Group(func(r chi.Router) {
		r.Get("/api/cities", s.GetCities)
		r.Get("/api/listings", s.GetListings)
		r.Post("/api/listings", s.CreateListing)
		r.Put("/api/listings/{id}", s.UpdateListing)
		r.Delete("/api/listings/{id}", s.DeleteListing)

		r.Get("/api/analytics", s.AnalyticsHandler)

	})
	r.With(s.AdminOnly).Get("/api/admin/users", s.AdminUsersHandler)
	r.With(s.AdminOnly).Get("/api/admin/listings", s.AdminListingsHandler)
	r.With(s.AdminOnly).Post("/api/admin/set-role", s.AdminSetRoleHandler)
	r.With(s.AdminOnly).Post("/api/admin/delete-user", s.AdminDeleteUserHandler)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
