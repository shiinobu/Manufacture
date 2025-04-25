package customer

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	EXE "id.benderaku.manufacture/app/helpers"
	customerModel "id.benderaku.manufacture/app/models/customer_model"
)

func ListCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customer, total, err := customerModel.GetListCustomer(w, r)
		if err != nil {
			EXE.ERROR.Println("Failed to fetch customers:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to fetch customers", "")
			return
		}
		EXE.SendList(w, http.StatusOK, total, "Fetched customers list", customer)
	}
}

func GetCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid customer ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid customer ID", "")
			return
		}
		user, err := customerModel.GetCustomerByID(id)
		if err != nil {
			EXE.ERROR.Println("Customer not found:", err)
			EXE.SendResponse(w, "", http.StatusNotFound, "Customer not found", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Fetched customer details", user)
	}
}

func CreateCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customer customerModel.Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		status, err := customerModel.CreateCustomer(&customer)
		if err != nil {
			EXE.ERROR.Println("Failed to create customer:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to create customer", "")
			return
		}
		if *status == 0 {
			EXE.ERROR.Println("Failed to create customer because code " + customer.Code + " already used")
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to create customer because code "+customer.Code+" already used", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Customer created successfully", customer)
	}
}

func UpdateCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid customer ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid customer ID", "")
			return
		}
		var customer customerModel.Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		if err := customerModel.UpdateCustomer(id, &customer); err != nil {
			EXE.ERROR.Println("Failed to update customer:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to update customer", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Customer updated successfully", customer)
	}
}

func DeleteCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid customer ID:", err)
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid customer ID", "")
			return
		}
		if err := customerModel.DeleteCustomer(id); err != nil {
			EXE.ERROR.Println("Failed to delete customer:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to delete customer", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Customer deleted successfully", "")
	}
}
