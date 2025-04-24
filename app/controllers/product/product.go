package product

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"id.benderaku.manufacture/app/helpers"
	productModel "id.benderaku.manufacture/app/models/product_model"
)

func ListProduct(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        pageStr := query.Get("page")
        limitStr := query.Get("limit")
        tipeStr := query.Get("type")
        page := 1
        limit := 10
        tipe := 0

        if pageStr != "" {
            pageInt, err := strconv.Atoi(pageStr)
            if err != nil || pageInt < 1 {
                errorLog.Println("Invalid page parameter:", pageStr)
                helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid page parameter", "")
                return
            }
            page = pageInt
        }

        if limitStr != "" {
            limitInt, err := strconv.Atoi(limitStr)
            if err != nil || limitInt < 1 {
                errorLog.Println("Invalid limit parameter:", limitStr)
                helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid limit parameter", "")
                return
            }
            limit = limitInt
        }

        if tipeStr != "" {
            tipeInt, err := strconv.Atoi(tipeStr)
            if err != nil || tipeInt < 0 {
                errorLog.Println("Invalid tipe parameter:", tipeStr)
                helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid tipe parameter", "")
                return
            }
            tipe = tipeInt
        }

        offset := (page - 1) * limit
        products, total, err := productModel.GetListProducts(tipe, limit, offset)
        if err != nil {
            errorLog.Println("Failed to fetch products: ", err)
            helpers.SendResponse(w, "", http.StatusInternalServerError, "Failed to fetch products", "")
            return
        }

        helpers.SendList(w, http.StatusOK, total, "Fetched products list", products)
    }
}

func GetProductMaterials(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        products, err := productModel.GetProductMaterials()
        if err != nil {
            errorLog.Println("Failed to getting products:", err)
            helpers.SendResponse(w, "", http.StatusInternalServerError, "Failed to getting products", "")
            return
        }
        w.Header().Set("Content-Type", "application/json")
        helpers.SendResponse(w, "", http.StatusOK, "Fetched products list", products)
    }
}

func GetProductById(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            errorLog.Println("Invalid product ID:", err)
            helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid product ID", "")
            return
        }
        product, err := productModel.GetProductByID(id)
        if err != nil {
            errorLog.Println("Product not found:", err)
            helpers.SendResponse(w, "", http.StatusNotFound, "Product not found", "")
            return
        }
        w.Header().Set("Content-Type", "application/json")
        helpers.SendResponse(w, "",  http.StatusOK, "Fetched product details", product)
    }
}

func CreateProduct(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var product productModel.Product
        if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
            errorLog.Println("Invalid request body:", err)
            helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
            return
        }
        if err := productModel.CreateProduct(&product); err != nil {
            errorLog.Println("Failed to create product:", err)
            helpers.SendResponse(w, "", http.StatusInternalServerError, "Failed to create product", "")
            return
        }
        helpers.SendResponse(w, "", http.StatusOK, "Product created successfully", product)
    }
}

func UpdateProduct(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            errorLog.Println("Invalid product ID:", err)
            http.Error(w, "Invalid product ID", http.StatusBadRequest)
            helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid product ID", "")
            return
        }
        var product productModel.Product
        if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
            errorLog.Println("Invalid request body:", err)
            helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid request body", "")
            return
        }
        if err := productModel.UpdateProduct(id, &product); err != nil {
            errorLog.Println("Failed to update product:", err)
            helpers.SendResponse(w, "", http.StatusInternalServerError, "Failed to update product", "")
            return
        }
        helpers.SendResponse(w, "", http.StatusOK, "Product updated successfully", product)
    }
}

func DeleteProduct(errorLog *log.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            errorLog.Println("Invalid product ID:", err)
            helpers.SendResponse(w, "", http.StatusBadRequest, "Invalid product ID", "")
            return
        }
        if err := productModel.DeleteProduct(id); err != nil {
            errorLog.Println("Failed to delete product:", err)
            helpers.SendResponse(w, "", http.StatusInternalServerError, "Failed to delete product", "")
            return
        }
        helpers.SendResponse(w, "", http.StatusOK, "Product deleted successfully", "")
    }
}