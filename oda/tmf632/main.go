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

type Individual struct {
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

// getIndividuals retrieves a individual
func getIndividuals(c *gin.Context) {
	var individuals []Individual

	query := `SELECT id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status FROM individual`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve individuals"})
		return
	}
	defer rows.Close()

	// Iterate over the result set and populate the slice
	for rows.Next() {
		var individual Individual
		if err := rows.Scan(&individual.ID, &individual.Gender, &individual.CountryOfBirth, &individual.Nationality, &individual.MaritalStatus, &individual.BirthDate, &individual.GivenName, &individual.PreferredGivenName, &individual.FamilyName, &individual.LegalName, &individual.MiddleName, &individual.FullName, &individual.FormattedName, &individual.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan individual"})
			return
		}

		// Set the custom field
		individual.Href = "http://localhost:8081/tmf-api/party/v5/individual/" + individual.ID
		individual.Type = "Individual"
		individual.BaseType = "Party"

		if errMsg := getExternalReference(&individual); errMsg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			return
		}

		// Append to the individuals slice
		individuals = append(individuals, individual)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching individuals"})
		return
	}

	c.JSON(http.StatusOK, individuals)
}

// getIndividualById retrieves a individual by ID
func getIndividualById(c *gin.Context) {
	id := c.Param("id")
	var individual Individual

	query := `SELECT id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status FROM individual WHERE id = $1 LIMIT 1`
	row := db.QueryRow(query, id)
	if err := row.Scan(&individual.ID, &individual.Gender, &individual.CountryOfBirth, &individual.Nationality, &individual.MaritalStatus, &individual.BirthDate, &individual.GivenName, &individual.PreferredGivenName, &individual.FamilyName, &individual.LegalName, &individual.MiddleName, &individual.FullName, &individual.FormattedName, &individual.Status); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Individual not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve individual"})
		return
	}

	// Set the custom field
	individual.Href = "http://localhost:8081/tmf-api/party/v5/individual/" + individual.ID
	individual.Type = "Individual"
	individual.BaseType = "Party"

	if errMsg := getExternalReference(&individual); errMsg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, individual)
}

// updateIndividual updates a individual
func updateIndividual(c *gin.Context) {
	// id := c.Param("id")
	var individual Individual

	// query := `SELECT id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status FROM individual WHERE id = $1 LIMIT 1`
	// row := db.QueryRow(query, id)
	// if err := row.Scan(&individual.ID, &individual.Gender, &individual.CountryOfBirth, &individual.Nationality, &individual.MaritalStatus, &individual.BirthDate, &individual.GivenName, &individual.PreferredGivenName, &individual.FamilyName, &individual.LegalName, &individual.MiddleName, &individual.FullName, &individual.FormattedName, &individual.Status); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "Individual not found"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve individual"})
	// 	return
	// }

	// // Set the custom field
	// individual.Href = "http://localhost:8081/tmf-api/party/v5/individual/" + individual.ID
	// individual.Type = "Individual"
	// individual.BaseType = "Party"

	// if errMsg := getExternalReference(&individual); errMsg != "" {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
	// 	return
	// }

	c.JSON(http.StatusOK, individual)
}

// createIndividual creates a new individual
func createIndividual(c *gin.Context) {
	var newIndividual Individual
	if err := c.ShouldBindJSON(&newIndividual); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO individual (id, gender, countryOfBirth, nationality, maritalStatus, birthDate, givenName, preferredGivenName, familyName, legalName, middleName, fullName, formattedName, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`
	err := db.QueryRow(query, newIndividual.ID, newIndividual.Gender, newIndividual.CountryOfBirth, newIndividual.Nationality, newIndividual.MaritalStatus, newIndividual.BirthDate, newIndividual.GivenName, newIndividual.PreferredGivenName, newIndividual.FamilyName, newIndividual.LegalName, newIndividual.MiddleName, newIndividual.FullName, newIndividual.FormattedName, newIndividual.Status).Scan(&newIndividual.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert individual" + err.Error()})
		return
	}

	if newIndividual.ExternalReferences != nil {
		for _, externalReference := range *newIndividual.ExternalReferences {
			var id int
			query = `INSERT INTO externalReference (name, externalIdentifierType, type, individual_id) VALUES ($1, $2, $3, $4) RETURNING id`
			err := db.QueryRow(query, externalReference.Name, externalReference.ExternalIdentifierType, externalReference.Type, newIndividual.ID).Scan(&id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert externalReference" + err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, newIndividual)
}

// deleteIndividualById deletes a individual by ID
func deleteIndividualById(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM externalReference WHERE individual_id = $1`
	res, err := db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count >= 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Delete externalReferences of individual: %s completed\n", id))
			}
		}
	}
	log.Println("no externalReferences of individual:", id)

	query = `DELETE FROM individual WHERE id = $1`
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
	c.JSON(http.StatusNotFound, gin.H{"error": "Individual not found"})
}

func getExternalReference(individual *Individual) string {
	var externalReferences []ExternalReference

	query := `SELECT name, externalIdentifierType, type FROM externalReference WHERE individual_id = $1`
	rows, err := db.Query(query, individual.ID)
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
		individual.ExternalReferences = &externalReferences
	}
	return ""
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/tmf-api/party/v5/individuals", getIndividuals)
	r.GET("/tmf-api/party/v5/individual/:id", getIndividualById)
	r.POST("/tmf-api/party/v5/individual", createIndividual)
	r.PATCH("/tmf-api/party/v5/individual", updateIndividual)
	r.DELETE("/tmf-api/party/v5/individual/:id", deleteIndividualById)
	r.Run(":8081")
}
