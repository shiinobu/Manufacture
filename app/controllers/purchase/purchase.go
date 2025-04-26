package purchase

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	EXE "id.benderaku.manufacture/app/helpers"
	purchaseModel "id.benderaku.manufacture/app/models/purchase_model"
)

func ListPurchase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		purchase, total, err := purchaseModel.GetListPurchase(w, r)
		if err != nil {
			EXE.ERROR.Println("Failed to fetch pruchases orders:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to fetch pruchases orders", "")
			return
		}
		EXE.SendList(w, http.StatusOK, total, "Fetched pruchases orders list", purchase)
	}
}

func GetPurchase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid purchase ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid purchase ID", "")
			return
		}
		purchase, err := purchaseModel.GetPurchaseByID(id)
		if err != nil {
			EXE.ERROR.Println("Supplier not found:", err)
			EXE.SendResponse(w, "", http.StatusNotFound, "Purchase orders not found", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Fetched purchase orders details", purchase)
	}
}

func CreatePurchase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var purchase purchaseModel.Purchase
		if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		status, err := purchaseModel.CreatePurchase(&purchase)
		if err != nil {
			EXE.ERROR.Println("Failed to create purchase:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to create purchase orders", "")
			return
		}
		if *status == 0 {
			EXE.ERROR.Println("Failed to create purchase orders because noreff " + purchase.Reff + " already used")
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to create purchase orders because noreff "+purchase.Reff+" already used", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Purchase orders created successfully", purchase)
	}
}
/*
func UpdatePurchase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid purchase ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid purchase ID", "")
			return
		}
		var purchase purchaseModel.Purchase
		if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		if err := purchaseModel.UpdatePurchase(id, &purchase); err != nil {
			EXE.ERROR.Println("Failed to update purchase orders:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to update purchase orders", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Purchase orders updated successfully", purchase)
	}
}

func DeletePurchase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid purchase ID:", err)
			http.Error(w, "Invalid purchase ID", http.StatusBadRequest)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid purchase ID", "")
			return
		}
		if err := purchaseModel.DeletePurchase(id); err != nil {
			EXE.ERROR.Println("Failed to delete purchase orders:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to delete purchase orders", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Purchase orders deleted successfully", "")
	}
}
*/
