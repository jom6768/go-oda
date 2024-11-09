package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jom6768/go-oda/oda/tmf632/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreateService creates a new service
func CreateService(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var service models.Service
		if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := db.Create(&service).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(service)
	}
}

// GetService retrieves a service by ID
func GetService(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var service models.Service

		if err := db.First(&service, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Service not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(service)
	}
}

// UpdateService updates an existing service
func UpdateService(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var service models.Service

		if err := db.First(&service, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Service not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db.Save(&service)
		json.NewEncoder(w).Encode(service)
	}
}

// DeleteService deletes a service by ID
func DeleteService(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if err := db.Delete(&models.Service{}, id).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
