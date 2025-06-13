package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"practic/internal/models"
	"time"
)

func NewToken(user models.UserDB, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Login
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
