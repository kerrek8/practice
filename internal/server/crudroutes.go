package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"practic/internal/logger/sl"
	"practic/internal/models"
	"strconv"
)

func (s *Server) CreateListing(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user").(*jwt.MapClaims)
	userIDparsed := *userID

	var l models.Listing
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		s.log.Error("Error in decoding body", sl.Err(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	l.UserID = int64(userIDparsed["uid"].(float64))

	uid, err := s.db.CreateListing(l.Name, l.Typel, l.Description, l.Status, l.City, l.Price, l.UserID)
	if err != nil {
		s.log.Error("Error in creating listing", sl.Err(err))
		http.Error(w, "Ошибка создания", 500)
		return
	}

	s.log.Info("Listing created successfully", slog.Int64("id", uid))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetListings(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	offset := int64((page - 1) * 10)
	filter := r.URL.Query().Get("filter")

	userID := r.Context().Value("user").(*jwt.MapClaims)
	userIDparsed := *userID
	userIDint := int64(userIDparsed["uid"].(float64))

	listings, err := s.db.GetListings(userIDint, offset, filter)

	if err != nil {
		s.log.Error("Error in getting listings", sl.Err(err))
		http.Error(w, "Ошибка получения списка", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listings); err != nil {
		s.log.Error("Error in encoding listings", sl.Err(err))
		http.Error(w, "Ошибка кодирования списка", 500)
		return
	}
}

func (s *Server) GetCities(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user").(*jwt.MapClaims)
	userIDparsed := *userID
	userIDint := int64(userIDparsed["uid"].(float64))

	var cities []string
	cities, err := s.db.GetCities(userIDint)
	if err != nil {
		s.log.Error("Error in getting cities", sl.Err(err))
		http.Error(w, "Ошибка получения городов", 500)
		return
	}

	json.NewEncoder(w).Encode(cities)
}

func (s *Server) UpdateListing(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user").(*jwt.MapClaims)
	userIDparsed := *userID

	var l models.Listing
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		s.log.Error("Error in decoding body", sl.Err(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	l.UserID = int64(userIDparsed["uid"].(float64))

	id := chi.URLParam(r, "id")
	listingID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		s.log.Error("Error in parsing listing ID", sl.Err(err))
		http.Error(w, "Invalid listing ID", http.StatusBadRequest)
		return
	}

	err = s.db.UpdateListing(l.Name, l.Typel, l.Description, l.Status, l.City, l.Price, listingID)
	if err != nil {
		s.log.Error("Error in updating listing", sl.Err(err))
		http.Error(w, "Ошибка обновления", 500)
		return
	}

	s.log.Info("Listing updated successfully", slog.Int64("id", listingID))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) DeleteListing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	listingID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		s.log.Error("Error in parsing listing ID", sl.Err(err))
		http.Error(w, "Invalid listing ID", http.StatusBadRequest)
		return
	}

	err = s.db.DeleteListing(listingID)
	if err != nil {
		s.log.Error("Error in deleting listing", sl.Err(err))
		http.Error(w, "Ошибка удаления", 500)
		return
	}

	s.log.Info("Listing deleted successfully", slog.Int64("id", listingID))
	w.WriteHeader(http.StatusNoContent)
}
