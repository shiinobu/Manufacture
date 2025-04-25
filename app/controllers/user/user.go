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
		users, total, err := userModel.GetListUsers(w, r)
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
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid user ID", "")
			return
		}
		user, err := userModel.GetUserByID(id)
		if err != nil {
			EXE.ERROR.Println("User not found:", err)
			EXE.SendResponse(w, "", http.StatusNotFound, "User not found", "")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		EXE.SendResponse(w, "", http.StatusOK, "Fetched users details", user)
	}
}

func CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user userModel.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		if err := userModel.CreateUser(&user); err != nil {
			EXE.ERROR.Println("Failed to create user:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to create user", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "User created successfully", user)
	}
}

func UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid user ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid user ID", "")
			return
		}
		var user userModel.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		if err := userModel.UpdateUser(id, &user); err != nil {
			EXE.ERROR.Println("Failed to update user:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to update user", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "User updated successfully", user)
	}
}

func DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid user ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid user ID", "")
			return
		}
		if err := userModel.DeleteUser(id); err != nil {
			EXE.ERROR.Println("Failed to delete user:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to delete user", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "User deleted successfully", "")
	}
}
