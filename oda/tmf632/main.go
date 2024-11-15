package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID                 string               `json:"id" binding:"required"`
	Href               string               `json:"href"`
	Type               string               `json:"@type"`
	BaseType           string               `json:"@baseType"`
	Gender             string               `json:"gender,omitempty"`
	CountryOfBirth     string               `json:"countryOfBirth,omitempty"`
	Nationality        string               `json:"nationality,omitempty"`
	MaritalStatus      string               `json:"maritalStatus,omitempty"`
	BirthDate          string               `json:"birthDate,omitempty"`
	GivenName          string               `json:"givenName,omitempty"`
	PreferredGivenName string               `json:"preferredGivenName,omitempty"`
	FamilyName         string               `json:"familyName,omitempty"`
	LegalName          string               `json:"legalName,omitempty"`
	MiddleName         string               `json:"middleName,omitempty"`
	FullName           string               `json:"fullName" binding:"required"`
	FormattedName      string               `json:"formattedName,omitempty"`
	Status             string               `json:"status,omitempty"`
	ExternalReferences *[]ExternalReference `json:"externalReference,omitempty"`
}

type ExternalReference struct {
	Name                   string `json:"id" binding:"required"`
	ExternalIdentifierType string `json:"externalIdentifierType"`
	Type                   string `json:"@type"`
}

var db *sql.DB

func initDB() {
	var err error
	// Run Local
	// connStr := "postgresql://myuser:mypass@localhost:5432/go_oda?sslmode=disable"
	// Run on Docker
	// connStr := "postgresql://myuser:mypass@host.docker.internal:5432/go_oda?sslmode=disable"
	// Run on Minikube
	connStr := "postgresql://myuser:mypass@host.minikube.internal:5432/go_oda?sslmode=disable"
	log.Println(connStr)
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			log.Println("Connected to the database successfully!")
			break
		}
		log.Println("Waiting for PostgreSQL to be ready...")
		time.Sleep(2 * time.Second)
	}
}

// createCustomer creates a new customer
func createCustomer(c *gin.Context) {
	var newCustomer Customer
	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO customer (id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`
	err := db.QueryRow(query, newCustomer.ID, newCustomer.Gender, newCustomer.CountryOfBirth, newCustomer.Nationality, newCustomer.MaritalStatus, newCustomer.BirthDate, newCustomer.GivenName, newCustomer.PreferredGivenName, newCustomer.FamilyName, newCustomer.LegalName, newCustomer.MiddleName, newCustomer.FullName, newCustomer.FormattedName, newCustomer.Status).Scan(&newCustomer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert customer" + err.Error()})
		return
	}

	if newCustomer.ExternalReferences != nil {
		for _, externalReference := range *newCustomer.ExternalReferences {
			var id int
			query = `INSERT INTO externalReference (name, externalIdentifierType, type, customer_id) VALUES ($1, $2, $3, $4) RETURNING id`
			err := db.QueryRow(query, externalReference.Name, externalReference.ExternalIdentifierType, externalReference.Type, newCustomer.ID).Scan(&id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert externalReference" + err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, newCustomer)
}

// getCustomers retrieves a customer
func getCustomers(c *gin.Context) {
	var customers []Customer

	query := `SELECT id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status FROM customer`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}
	defer rows.Close()

	// Iterate over the result set and populate the slice
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Gender, &customer.CountryOfBirth, &customer.Nationality, &customer.MaritalStatus, &customer.BirthDate, &customer.GivenName, &customer.PreferredGivenName, &customer.FamilyName, &customer.LegalName, &customer.MiddleName, &customer.FullName, &customer.FormattedName, &customer.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan customer"})
			return
		}

		// Set the custom field
		customer.Href = "https://serverRoot/tmf-api/party/v5/individual/" + customer.ID
		customer.Type = "Individual"
		customer.BaseType = "Party"

		if errMsg := getExternalReference(&customer); errMsg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			return
		}

		// Append to the customers slice
		customers = append(customers, customer)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching customers"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// getCustomerById retrieves a customer by ID
func getCustomerById(c *gin.Context) {
	id := c.Param("id")
	var customer Customer

	query := `SELECT id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status FROM customer WHERE id = $1 LIMIT 1`
	row := db.QueryRow(query, id)
	if err := row.Scan(&customer.ID, &customer.Gender, &customer.CountryOfBirth, &customer.Nationality, &customer.MaritalStatus, &customer.BirthDate, &customer.GivenName, &customer.PreferredGivenName, &customer.FamilyName, &customer.LegalName, &customer.MiddleName, &customer.FullName, &customer.FormattedName, &customer.Status); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer"})
		return
	}

	// Set the custom field
	customer.Href = "https://serverRoot/tmf-api/party/v5/individual/" + customer.ID
	customer.Type = "Individual"
	customer.BaseType = "Party"

	if errMsg := getExternalReference(&customer); errMsg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// deleteCustomerById deletes a customer by ID
func deleteCustomerById(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM externalReference WHERE customer_id = $1`
	res, err := db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count >= 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Delete externalReferences of customer: %s completed\n", id))
			}
		}
	}
	log.Println("no externalReferences of customer:", id)

	query = `DELETE FROM customer WHERE id = $1`
	res, err = db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count == 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("id = %s completed.", id))
				return
			}
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
}

func getExternalReference(customer *Customer) string {
	var externalReferences []ExternalReference

	query := `SELECT name, externalIdentifierType, type FROM externalReference WHERE customer_id = $1`
	rows, err := db.Query(query, customer.ID)
	if err != nil {
		return "Failed to retrieve externalReferences"
	}
	defer rows.Close()

	// Iterate over the result set and populate the slice
	for rows.Next() {
		var externalReference ExternalReference
		if err := rows.Scan(&externalReference.Name, &externalReference.ExternalIdentifierType, &externalReference.Type); err != nil {
			return "Failed to scan externalReference"
		}

		// Append to the externalReferences slice
		externalReferences = append(externalReferences, externalReference)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return "Error while fetching externalReferences"
	}

	if externalReferences != nil {
		customer.ExternalReferences = &externalReferences
	}
	return ""
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/tmf632/customer", createCustomer)
	r.GET("/tmf632/customers", getCustomers)
	r.GET("/tmf632/customer/:id", getCustomerById)
	r.DELETE("/tmf632/customer/:id", deleteCustomerById)
	r.Run(":8081")
}
