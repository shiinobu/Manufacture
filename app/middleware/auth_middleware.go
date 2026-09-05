package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"id.benderaku.manufacture/app/helpers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("token")
		if authHeader == "" {
			helpers.SendResponse(w, "", http.StatusUnauthorized, "Missing Token header", "")
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.SendResponse(w, "", http.StatusUnauthorized, "Invalid Token header format", "")
			return
		}

		token, err := helpers.ParseToken(parts[1])
		if err != nil || !token.Valid {
			helpers.SendResponse(w, "", http.StatusUnauthorized, "Invalid or expired token", "")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			helpers.SendResponse(w, "", http.StatusUnauthorized, "Invalid token claims", "")
			return
		}

		userID, ok := claims["sub"].(float64)
		if !ok {
			helpers.SendResponse(w, "", http.StatusUnauthorized, "Invalid user ID in token", "")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", int(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
