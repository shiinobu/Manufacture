package supplier

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	EXE "id.benderaku.manufacture/app/helpers"
	supplierModel "id.benderaku.manufacture/app/models/supplier_model"
)

func ListSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supplier, total, err := supplierModel.GetListSupplier(w, r)
		if err != nil {
			EXE.ERROR.Println("Failed to fetch suppliers:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to fetch suppliers", "")
			return
		}
		EXE.SendList(w, http.StatusOK, total, "Fetched suppliers list", supplier)
	}
}

func GetSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid supplier ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid supplier ID", "")
			return
		}
		user, err := supplierModel.GetSupplierByID(id)
		if err != nil {
			EXE.ERROR.Println("Supplier not found:", err)
			EXE.SendResponse(w, "", http.StatusNotFound, "Supplier not found", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Fetched supplier details", user)
	}
}

func CreateSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var supplier supplierModel.Supplier
		if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		if err := supplierModel.CreateSupplier(&supplier); err != nil {
			EXE.ERROR.Println("Failed to create supplier:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to create supplier", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Supplier created successfully", supplier)
	}
}

func UpdateSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid supplier ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid supplier ID", "")
			return
		}
		var supplier supplierModel.Supplier
		if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		if err := supplierModel.UpdateSupplier(id, &supplier); err != nil {
			EXE.ERROR.Println("Failed to update supplier:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to update supplier", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Supplier updated successfully", supplier)
	}
}

func DeleteSupplier() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid supplier ID:", err)
			http.Error(w, "Invalid supplier ID", http.StatusBadRequest)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid supplier ID", "")
			return
		}
		if err := supplierModel.DeleteSupplier(id); err != nil {
			EXE.ERROR.Println("Failed to delete supplier:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to delete supplier", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Supplier deleted successfully", "")
	}
}
