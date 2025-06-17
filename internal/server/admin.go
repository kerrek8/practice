package server

import (
	"encoding/json"
	"net/http"
	"practic/internal/logger/sl"
)

func (s *Server) AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := s.db.GetAllUsers()
	if err != nil {
		s.log.Error("Error fetching users", sl.Err(err))
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(users)
	if err != nil {
		s.log.Error("Error marshalling users", sl.Err(err))
		http.Error(w, "Failed to process users data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)

}

func (s *Server) AdminListingsHandler(w http.ResponseWriter, r *http.Request) {
	listings, err := s.db.GetAllListings()
	if err != nil {
		s.log.Error("Error fetching listings", sl.Err(err))
		http.Error(w, "Failed to fetch listings", http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(listings)
	if err != nil {
		s.log.Error("Error marshalling listings", sl.Err(err))
		http.Error(w, "Failed to process listings data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)
}

func (s *Server) AdminSetRoleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64  `json:"user_id"`
		Role   string `json:"role"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.log.Error("Error decoding request body", sl.Err(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = s.db.SetUserRole(req.UserID, req.Role)
	if err != nil {
		s.log.Error("Error setting user role", sl.Err(err))
		http.Error(w, "Failed to set user role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) AdminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int64 `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.log.Error("Error decoding request body", sl.Err(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = s.db.DeleteUser(req.UserID)
	if err != nil {
		s.log.Error("Error deleting user", sl.Err(err))
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
