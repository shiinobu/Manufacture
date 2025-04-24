package helpers

import (
	"time"

	"id.benderaku.manufacture/app/config"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte(config.JWT_KEY) // Store securely in environment variables

func GenerateToken(userID int) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(JwtSecret)
}