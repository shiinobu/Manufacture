package helpers

import (
	"testing"
	"time"

	"id.benderaku.manufacture/app/config"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseToken(t *testing.T) {
	t.Setenv("JWT_KEY", "test-secret")

	tokenString, err := GenerateToken(123)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	token, err := ParseToken(tokenString)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if !token.Valid {
		t.Fatal("expected token to be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected jwt.MapClaims")
	}

	if got, ok := claims["sub"].(float64); !ok || int(got) != 123 {
		t.Fatalf("unexpected subject claim: %#v", claims["sub"])
	}
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	t.Setenv("JWT_KEY", "correct-secret")

	tokenString, err := GenerateToken(123)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	t.Setenv("JWT_KEY", "wrong-secret")

	token, err := ParseToken(tokenString)
	if err == nil || token.Valid {
		t.Fatal("expected token signed with a different secret to be rejected")
	}
}

func TestParseTokenRejectsExpiredToken(t *testing.T) {
	t.Setenv("JWT_KEY", "test-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 123,
		"exp": time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(config.JWTKey())
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	parsed, err := ParseToken(tokenString)
	if err == nil || parsed.Valid {
		t.Fatal("expected expired token to be rejected")
	}
}

func TestParseTokenRejectsUnexpectedSigningMethod(t *testing.T) {
	t.Setenv("JWT_KEY", "test-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{
		"sub": 123,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(config.JWTKey())
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	parsed, err := ParseToken(tokenString)
	if err == nil || parsed.Valid {
		t.Fatal("expected unexpected signing method to be rejected")
	}
}
