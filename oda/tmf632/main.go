package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jom6768/go-oda/oda/tmf632/config"
	"github.com/jom6768/go-oda/oda/tmf632/handlers"
	"github.com/jom6768/go-oda/oda/tmf632/models"

	"github.com/gorilla/mux"
)

func main() {
	db := config.ConnectDB()
	models.MigrateService(db) // Run database migration

	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/tmf632/services", handlers.CreateService(db)).Methods("POST")
	router.HandleFunc("/tmf632/services/{id}", handlers.GetService(db)).Methods("GET")
	router.HandleFunc("/tmf632/services/{id}", handlers.UpdateService(db)).Methods("PUT")
	router.HandleFunc("/tmf632/services/{id}", handlers.DeleteService(db)).Methods("DELETE")

	// Start server
	fmt.Println("TMF632 Service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
