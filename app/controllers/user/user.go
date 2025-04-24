package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	EXE "id.benderaku.manufacture/app/helpers"
	userModel "id.benderaku.manufacture/app/models/user_model"
)

func ListUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, total, err := userModel.GetAllUsers(w, r)
		if err != nil {
			EXE.ERROR.Println("Failed to fetch users: ", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to fetch users", "")
			return
		}

		EXE.SendList(w, http.StatusOK, total, "Fetched users list", users)
	}
}

func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid user ID:", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		user, err := userModel.GetUserByID(id)
		if err != nil {
			EXE.ERROR.Println("User not found:", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user userModel.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		isActive, err := userModel.CheckStatusUser(user.Username, user.Email)
		if err != nil {
			if err := userModel.CreateUser(&user); err != nil {
				EXE.ERROR.Println("Failed to create user:", err)
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
		} else {
			if isActive.Status == 0 {
				if err := userModel.CreateUser(&user); err != nil {
					EXE.ERROR.Println("Failed to create user:", err)
					http.Error(w, "Failed to create user", http.StatusInternalServerError)
					return
				}
			} else {
				if err := userModel.UpdateStatusUser(isActive.ID, &user); err != nil {
					EXE.ERROR.Println("Failed to create user:", err)
					http.Error(w, "Failed to create user", http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid user ID:", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		var user userModel.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := userModel.UpdateUser(id, &user); err != nil {
			EXE.ERROR.Println("Failed to update user:", err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid user ID:", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		if err := userModel.DeleteUser(id); err != nil {
			EXE.ERROR.Println("Failed to delete user:", err)
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
