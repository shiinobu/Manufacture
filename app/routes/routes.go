package routes

import (
	"github.com/gorilla/mux"

	Auth "id.benderaku.manufacture/app/controllers/auth"
	Customer "id.benderaku.manufacture/app/controllers/customer"
	Product "id.benderaku.manufacture/app/controllers/product"
	Supplier "id.benderaku.manufacture/app/controllers/supplier"
	User "id.benderaku.manufacture/app/controllers/user"
	Protector "id.benderaku.manufacture/app/middleware"
)

var Router *mux.Router

func RegisterRoutes() {
	Router = mux.NewRouter()
	publicRouter := Router.PathPrefix("/").Subrouter()
	protectedRouter := Router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(Protector.AuthMiddleware)

	// API ROUTES FOR LOGIN AND LOGOUT
	publicRouter.HandleFunc("/login", Auth.Login()).Methods("POST")
	protectedRouter.HandleFunc("/logout", Auth.Logout()).Methods("POST")

	// API ROUTES FOR USERS
	protectedRouter.HandleFunc("/users", User.ListUser()).Methods("GET")
	publicRouter.HandleFunc("/users", User.CreateUser()).Methods("POST")
	protectedRouter.HandleFunc("/users/{id}", User.GetUser()).Methods("GET")
	protectedRouter.HandleFunc("/users/{id}", User.UpdateUser()).Methods("PUT")
	protectedRouter.HandleFunc("/users/{id}", User.DeleteUser()).Methods("DELETE")

	// API ROUTES FOR PRODUCTS
	protectedRouter.HandleFunc("/products", Product.ListProduct()).Methods("GET")
	protectedRouter.HandleFunc("/products", Product.CreateProduct()).Methods("POST")
	protectedRouter.HandleFunc("/products/materials", Product.GetProductMaterials()).Methods("GET")
	protectedRouter.HandleFunc("/products/{id}", Product.GetProductById()).Methods("GET")
	protectedRouter.HandleFunc("/products/{id}", Product.UpdateProduct()).Methods("PUT")
	protectedRouter.HandleFunc("/products/{id}", Product.DeleteProduct()).Methods("DELETE")

	// API ROUTES FOR SUPPLIERS
	protectedRouter.HandleFunc("/suppliers", Supplier.ListSupplier()).Methods("GET")
	protectedRouter.HandleFunc("/suppliers", Supplier.CreateSupplier()).Methods("POST")
	protectedRouter.HandleFunc("/suppliers/{id}", Supplier.GetSupplier()).Methods("GET")
	protectedRouter.HandleFunc("/suppliers/{id}", Supplier.UpdateSupplier()).Methods("PUT")
	protectedRouter.HandleFunc("/suppliers/{id}", Supplier.DeleteSupplier()).Methods("DELETE")

	// API ROUTES FOR CUSTOMERS
	protectedRouter.HandleFunc("/customers", Customer.ListCustomer()).Methods("GET")
	protectedRouter.HandleFunc("/customers", Customer.CreateCustomer()).Methods("POST")
	protectedRouter.HandleFunc("/customers/{id}", Customer.GetCustomer()).Methods("GET")
	protectedRouter.HandleFunc("/customers/{id}", Customer.UpdateCustomer()).Methods("PUT")
	protectedRouter.HandleFunc("/customers/{id}", Customer.DeleteCustomer()).Methods("DELETE")
}
