package server

import (
	"encoding/json"
	"net/http"
	"practic/internal/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	r.Post("/api/logout", LogoutHandler)
	//
	//r.Group(func(r chi.Router) {
	//
	//	r.Get("/api/listings", GetListings)
	//	r.Post("/api/listings", CreateListing)
	//	r.Put("/api/listings/{id}", UpdateListing)
	//	r.Delete("/api/listings/{id}", DeleteListing)
	//
	//	r.Get("/api/analytics", AnalyticsHandler)
	//})

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		s.log.Error("error handling JSON marshal. Err: ", sl.Err(err))
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
