package server

import (
	"encoding/json"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"practic/internal/jwt"
	"practic/internal/logger/sl"
	"practic/internal/models"
	"time"
)

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		s.log.Error("Error in decoding body", sl.Err(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("Error in hashing password", sl.Err(err))
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
	}
	uid, err := s.db.CreateUser(u.Name, u.Login, passHash)
	if err != nil {
		s.log.Error("Error in creating user", sl.Err(err))
		http.Error(w, fmt.Sprintf("failed to create user: %v", err), http.StatusInternalServerError)
	}
	s.log.Info("User created successfully", slog.Int64("id", uid))
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		s.log.Error("Error in decoding body", sl.Err(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := s.db.User(u.Login)
	if err != nil {
		s.log.Error("Error in getting user", sl.Err(err))
		http.Error(w, fmt.Sprintf("failed to get user: %v", err), http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		s.log.Error("Error in comparing password", sl.Err(err))
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := jwt.NewToken(user, time.Hour)
	if err != nil {
		s.log.Error("Error in creating token", sl.Err(err))
		http.Error(w, fmt.Sprintf("failed to create token: %v", err), http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:  "token",
		Path:  "/",
		Value: token}
	http.SetCookie(w, cookie)
	_, err = w.Write([]byte("Login successful"))
	if err != nil {
		return
	}
	s.log.Info("User logged in", slog.Int64("id", user.ID))
	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {}
