package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	userModel "id.benderaku.manufacture/app/models/user_model"
)

func ListUser(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        users, err := userModel.GetAllUsers()
        if err != nil {
            errorLog.Println("Failed to fetch users:", err)
            http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}

func GetUser(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            errorLog.Println("Invalid user ID:", err)
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }
        user, err := userModel.GetUserByID(id)
        if err != nil {
            errorLog.Println("User not found:", err)
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}

func CreateUser(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user userModel.User

        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            errorLog.Println("Invalid request body:", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        isActive, err := userModel.CheckStatusUser(user.Username, user.Email)
        if err != nil {
            if err := userModel.CreateUser(&user); err != nil {
                errorLog.Println("Failed to create user:", err)
                http.Error(w, "Failed to create user", http.StatusInternalServerError)
                return
            }
        } else {
            if isActive.Status == 0 {
                if err := userModel.CreateUser(&user); err != nil {
                    errorLog.Println("Failed to create user:", err)
                    http.Error(w, "Failed to create user", http.StatusInternalServerError)
                    return
                }
            } else {
                if err := userModel.UpdateStatusUser(isActive.ID, &user); err != nil {
                    errorLog.Println("Failed to create user:", err)
                    http.Error(w, "Failed to create user", http.StatusInternalServerError)
                    return
                }
            }
        }

        w.WriteHeader(http.StatusCreated)
    }
}

func UpdateUser(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            errorLog.Println("Invalid user ID:", err)
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }
        var user userModel.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            errorLog.Println("Invalid request body:", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        if err := userModel.UpdateUser(id, &user); err != nil {
            errorLog.Println("Failed to update user:", err)
            http.Error(w, "Failed to update user", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
    }
}

func DeleteUser(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            errorLog.Println("Invalid user ID:", err)
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }
        if err := userModel.DeleteUser(id); err != nil {
            errorLog.Println("Failed to delete user:", err)
            http.Error(w, "Failed to delete user", http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusNoContent)
    }
}