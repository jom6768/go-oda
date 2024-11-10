package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jom6768/go-oda/oda/tmf669/config"
	"github.com/jom6768/go-oda/oda/tmf669/handlers"
	"github.com/jom6768/go-oda/oda/tmf669/models"

	"github.com/gorilla/mux"
)

func main() {
	db := config.ConnectDB()
	models.MigrateProduct(db) // Run database migration

	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/tmf669/products", handlers.CreateProduct(db)).Methods("POST")
	router.HandleFunc("/tmf669/products/{id}", handlers.GetProduct(db)).Methods("GET")
	router.HandleFunc("/tmf669/products/{id}", handlers.UpdateProduct(db)).Methods("PUT")
	router.HandleFunc("/tmf669/products/{id}", handlers.DeleteProduct(db)).Methods("DELETE")

	// Start server
	fmt.Println("TMF669 Service running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
