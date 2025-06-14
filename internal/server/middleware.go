package server

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"os"
	"practic/internal/logger/sl"
)

var jwtKey = os.Getenv("JWT_KEY")

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/register", "/", "/api/register", "/api/login", "/styles.css", "/app.js", "/register/", "/api/logout", "/health"}
		requestPath := r.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			s.log.Error("Error in getting cookie", sl.Err(err))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if err != nil || !token.Valid {
			s.log.Error("Error in parsing token", sl.Err(err), slog.Any("token", token))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
