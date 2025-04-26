package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	EXE "id.benderaku.manufacture/app/helpers"
	userModel "id.benderaku.manufacture/app/models/user_model"
)

func Activation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		user, err := userModel.GetUserByEmail(creds.Email)
		if err != nil {
			EXE.ERROR.Println("User not found:", err)
			EXE.SendResponse(w, "", http.StatusUnauthorized, "User not found", "")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			EXE.ERROR.Println("Invalid password:", err)
			EXE.SendResponse(w, "", http.StatusUnauthorized, "Invalid password", "")
			return
		}
		token, err := EXE.GenerateToken(user.ID)
		if err != nil {
			EXE.ERROR.Println("Failed to generate token:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to generated token", "")
			return
		}
		EXE.SendResponse(w, token, http.StatusOK, "Logged in successfully", "")
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var set string
		var raw map[string]any
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			EXE.ERROR.Println("Failed to read request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Failed to read request body", "")
			return
		}
		defer r.Body.Close()
		if len(bodyBytes) == 0 {
			EXE.ERROR.Println("Empty request body")
			EXE.SendResponse(w, "", http.StatusBadRequest, "Empty request body", "")
			return
		}
		if err := json.Unmarshal(bodyBytes, &raw); err != nil {
			EXE.ERROR.Println("Invalid JSON format:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid JSON format", "")
			return
		}
		found := false
		for key := range raw {
			if key == "email" {
				found = true
				set = "email"
				break
			}
			if key == "password" {
				found = true
				set = "password"
				break
			}
		}
		if !found {
			EXE.ERROR.Println("JSON key " + set + " not found or incorrect")
			EXE.SendResponse(w, "", http.StatusBadRequest, "Expected JSON key " + set, "")
			return
		}
		if err := json.Unmarshal(bodyBytes, &creds); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		creds.Email = strings.TrimSpace(creds.Email)
		creds.Email = strings.ToLower(creds.Email)
		creds.Password = strings.TrimSpace(creds.Password)
		user, err := userModel.GetUserByEmail(creds.Email)
		if err != nil {
			EXE.ERROR.Println("User not found:", err)
			EXE.SendResponse(w, "", http.StatusUnauthorized, "User not found", "")
			return
		}
		if user.Islogin == 0 {
			EXE.ERROR.Println("User not activated, please activated first")
			EXE.SendResponse(w, "", http.StatusUnauthorized, "User not activated, please activated first", "")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			EXE.ERROR.Println("Invalid password:", err)
			EXE.SendResponse(w, "", http.StatusUnauthorized, "Invalid password", "")
			return
		}
		token, err := EXE.GenerateToken(user.ID)
		if err != nil {
			EXE.ERROR.Println("Failed to generate token:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to generated token", "")
			return
		}
		EXE.SendResponse(w, token, http.StatusOK, "Logged in successfully", "")
	}
}

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID")
		if userID == nil {
			EXE.SendResponse(w, "", http.StatusInternalServerError, "User ID not found in context", "")
		}
		w.Header().Set("Content-Type", "application/json")
		EXE.SendResponse(w, "", http.StatusOK, "Logged out successfully", "")
	}
}
