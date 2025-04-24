package supplier

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	Log "id.benderaku.manufacture/app/helpers"
	// Resp "id.benderaku.manufacture/app/helpers"
	supplierModel "id.benderaku.manufacture/app/models/supplier_model"
)

func ListSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := supplierModel.GetAllSupplier()
		if err != nil {
			Log.ERROR.Println("Failed to fetch suppliers:", err)
			http.Error(w, "Failed to fetch suppliers", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			Log.ERROR.Println("Invalid supplier ID:", err)
			http.Error(w, "Invalid supplier ID", http.StatusBadRequest)
			return
		}
		user, err := supplierModel.GetSupplierByID(id)
		if err != nil {
			Log.ERROR.Println("Supplier not found:", err)
			http.Error(w, "Supplier not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func CreateSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user supplierModel.Supplier

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			Log.ERROR.Println("Invalid request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		isActive, err := supplierModel.CheckStatusSupplier(user.Username, user.Email)
		if err != nil {
			if err := supplierModel.CreateSupplier(&user); err != nil {
				Log.ERROR.Println("Failed to create supplier:", err)
				http.Error(w, "Failed to create supplier", http.StatusInternalServerError)
				return
			}
		} else {
			if isActive.Status == 0 {
				if err := supplierModel.CreateSupplier(&user); err != nil {
					Log.ERROR.Println("Failed to create supplier:", err)
					http.Error(w, "Failed to create supplier", http.StatusInternalServerError)
					return
				}
			} else {
				if err := supplierModel.UpdateStatusSupplier(isActive.ID, &user); err != nil {
					Log.ERROR.Println("Failed to update supplier status:", err)
					http.Error(w, "Failed to update supplier status", http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			Log.ERROR.Println("Invalid user ID:", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		var user supplierModel.Supplier
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			Log.ERROR.Println("Invalid request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := supplierModel.UpdateSupplier(id, &user); err != nil {
			Log.ERROR.Println("Failed to update supplier:", err)
			http.Error(w, "Failed to update supplier", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			Log.ERROR.Println("Invalid supplier ID:", err)
			http.Error(w, "Invalid supplier ID", http.StatusBadRequest)
			return
		}
		if err := supplierModel.DeleteSupplier(id); err != nil {
			Log.ERROR.Println("Failed to delete supplier:", err)
			http.Error(w, "Failed to delete supplier", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
