package routes

import (
	"log"

	"github.com/gorilla/mux"

	Auth "id.benderaku.manufacture/app/controllers/auth"
	Product "id.benderaku.manufacture/app/controllers/product"
	Supplier "id.benderaku.manufacture/app/controllers/supplier"
	User "id.benderaku.manufacture/app/controllers/user"
	Protector "id.benderaku.manufacture/app/middleware"
)

// var Router *mux.Router

func RegisterRoutes(r *mux.Router, infoLog, errorLog *log.Logger) {
    publicRouter := r.PathPrefix("/").Subrouter()
    protectedRouter := r.PathPrefix("/api").Subrouter()
    protectedRouter.Use(Protector.AuthMiddleware)

    // API ROUTES FOR LOGIN AND LOGOUT
    publicRouter.HandleFunc("/login", Auth.Login(errorLog)).Methods("POST")
    protectedRouter.HandleFunc("/logout", Auth.Logout(errorLog)).Methods("POST")

    // API ROUTES FOR USERS
    protectedRouter.HandleFunc("/users", User.ListUser(errorLog)).Methods("GET")
    publicRouter.HandleFunc("/users", User.CreateUser(errorLog)).Methods("POST")
    protectedRouter.HandleFunc("/users/{id}", User.GetUser(errorLog)).Methods("GET")
    protectedRouter.HandleFunc("/users/{id}", User.UpdateUser(errorLog)).Methods("PUT")
    protectedRouter.HandleFunc("/users/{id}", User.DeleteUser(errorLog)).Methods("DELETE")
    
    // API ROUTES FOR PRODUCTS
    protectedRouter.HandleFunc("/products", Product.ListProduct(errorLog)).Methods("GET")
    protectedRouter.HandleFunc("/products", Product.CreateProduct(errorLog)).Methods("POST")
    protectedRouter.HandleFunc("/products/{id}", Product.GetProductById(errorLog)).Methods("GET")
    protectedRouter.HandleFunc("/products/{id}", Product.UpdateProduct(errorLog)).Methods("PUT")
    protectedRouter.HandleFunc("/products/{id}", Product.DeleteProduct(errorLog)).Methods("DELETE")
    protectedRouter.HandleFunc("/products/materials", Product.GetProductMaterials(errorLog)).Methods("POST")
    
    // API ROUTES FOR SUPPLIERS
    protectedRouter.HandleFunc("/suppliers", Supplier.ListSupplier(errorLog)).Methods("GET")
    protectedRouter.HandleFunc("/suppliers", Supplier.CreateSupplier(errorLog)).Methods("POST")
    protectedRouter.HandleFunc("/suppliers/{id}", Supplier.GetSupplier(errorLog)).Methods("GET")
    protectedRouter.HandleFunc("/suppliers/{id}", Supplier.UpdateSupplier(errorLog)).Methods("PUT")
    protectedRouter.HandleFunc("/suppliers/{id}", Supplier.DeleteSupplier(errorLog)).Methods("DELETE")
    
}