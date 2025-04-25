package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	EXE "id.benderaku.manufacture/app/helpers"
	productModel "id.benderaku.manufacture/app/models/product_model"
)

func ListProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, total, err := productModel.GetListProducts(w, r)
		if err != nil {
			EXE.ERROR.Println("Failed to fetch products: ", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to fetch products", "")
			return
		}
		EXE.SendList(w, http.StatusOK, total, "Fetched products list", products)
	}
}

func GetProductMaterials() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := productModel.GetProductMaterials()
		if err != nil {
			EXE.ERROR.Println("Failed to getting products:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to getting products", "")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		EXE.SendResponse(w, "", http.StatusOK, "Fetched products list", products)
	}
}

func GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid product ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid product ID", "")
			return
		}
		product, err := productModel.GetProductByID(id)
		if err != nil {
			EXE.ERROR.Println("Product not found:", err)
			EXE.SendResponse(w, "", http.StatusNotFound, "Product not found", "")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		EXE.SendResponse(w, "", http.StatusOK, "Fetched product details", product)
	}
}

func CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product productModel.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		if err := productModel.CreateProduct(&product); err != nil {
			EXE.ERROR.Println("Failed to create product:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to create product", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Product created successfully", product)
	}
}

func UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid product ID:", err)
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid product ID", "")
			return
		}
		var product productModel.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			EXE.ERROR.Println("Invalid request body:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
			return
		}
		if err := productModel.UpdateProduct(id, &product); err != nil {
			EXE.ERROR.Println("Failed to update product:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to update product", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Product updated successfully", product)
	}
}

func DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			EXE.ERROR.Println("Invalid product ID:", err)
			EXE.SendResponse(w, "", http.StatusBadRequest, "Invalid product ID", "")
			return
		}
		if err := productModel.DeleteProduct(id); err != nil {
			EXE.ERROR.Println("Failed to delete product:", err)
			EXE.SendResponse(w, "", http.StatusInternalServerError, "Failed to delete product", "")
			return
		}
		EXE.SendResponse(w, "", http.StatusOK, "Product deleted successfully", "")
	}
}
