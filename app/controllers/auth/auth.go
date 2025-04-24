package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"id.benderaku.manufacture/app/helpers"
	userModel "id.benderaku.manufacture/app/models/user_model"
)

func Activation(db *sql.DB, infoLog, errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var creds struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
            errorLog.Println("Invalid request body:", err)
            helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
            return
        }

        user, err := userModel.GetUserByUsername(creds.Username)
        if err != nil {
            errorLog.Println("User not found:", err)
            helpers.SendResponse(w, "", http.StatusUnauthorized, "User not found", "")
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
            errorLog.Println("Invalid password:", err)
            helpers.SendResponse(w, "", http.StatusUnauthorized, "Invalid password", "")
            return
        }

        token, err := helpers.GenerateToken(user.ID)
        if err != nil {
            errorLog.Println("Failed to generate token:", err)
            helpers.SendResponse(w, "", http.StatusInternalServerError, "Failed to generated token", "")
            return
        }

        helpers.SendResponse(w, token, http.StatusOK, "Logged in successfully", "")
    }
}

func Login(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var creds struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
            errorLog.Println("Invalid request body:", err)
            helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
            return
        }

        user, err := userModel.GetUserByUsername(creds.Username)
        if err != nil {
            errorLog.Println("User not found:", err)
            helpers.SendResponse(w, "", http.StatusUnauthorized, "User not found", "")
            return
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
            errorLog.Println("Invalid password:", err)
            helpers.SendResponse(w, "", http.StatusUnauthorized, "Invalid password", "")
            return
        }

        token, err := helpers.GenerateToken(user.ID)
        if err != nil {
            errorLog.Println("Failed to generate token:", err)
            helpers.SendResponse(w, "", http.StatusInternalServerError, "Failed to generated token", "")
            return
        }

        helpers.SendResponse(w, token, http.StatusOK, "Logged in successfully", "")
    }
}

func Logout(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.Context().Value("userID")
        if userID == nil {
            helpers.SendResponse(w, "", http.StatusInternalServerError, "User ID not found in context", "")
        }
        w.Header().Set("Content-Type", "application/json")
        helpers.SendResponse(w, "", http.StatusOK, "Logged out successfully", "")
    }
}