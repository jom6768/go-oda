package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Individual struct {
	ID                 *string              `json:"id" binding:"required"`
	Href               *string              `json:"href,omitempty"`
	Type               *string              `json:"@type" binding:"required"`
	BaseType           *string              `json:"@baseType,omitempty"`
	Gender             *string              `json:"gender,omitempty"`
	CountryOfBirth     *string              `json:"countryOfBirth,omitempty"`
	Nationality        *string              `json:"nationality,omitempty"`
	MaritalStatus      *string              `json:"maritalStatus,omitempty"`
	BirthDate          *string              `json:"birthDate,omitempty"`
	GivenName          *string              `json:"givenName,omitempty"`
	PreferredGivenName *string              `json:"preferredGivenName,omitempty"`
	FamilyName         *string              `json:"familyName,omitempty"`
	LegalName          *string              `json:"legalName,omitempty"`
	MiddleName         *string              `json:"middleName,omitempty"`
	FullName           *string              `json:"fullName,omitempty"`
	FormattedName      *string              `json:"formattedName,omitempty"`
	Status             *string              `json:"status,omitempty"` //"initialized","validated","deceaded"
	ExternalReferences *[]ExternalReference `json:"externalReference,omitempty"`
}

type Organization struct {
	ID                 *string              `json:"id" binding:"required"`
	Href               *string              `json:"href,omitempty"`
	Type               *string              `json:"@type" binding:"required"`
	BaseType           *string              `json:"@baseType,omitempty"`
	IsLegalEntity      *bool                `json:"isLegalEntity,omitempty"`
	IsHeadOffice       *bool                `json:"isHeadOffice,omitempty"`
	OrganizationType   *string              `json:"organizationType,omitempty"`
	Name               *string              `json:"name,omitempty"`
	TradingName        *string              `json:"tradingName,omitempty"`
	NameType           *string              `json:"nameType,omitempty"`
	Status             *string              `json:"status,omitempty"` //"initialized","validated","closed"
	ExternalReferences *[]ExternalReference `json:"externalReference,omitempty"`
}

type ExternalReference struct {
	Name                   string `json:"name" binding:"required"`
	ExternalIdentifierType string `json:"externalIdentifierType,omitempty"`
	Type                   string `json:"@type,omitempty"`
}

var db *sql.DB

func initDB() {
	var err error
	// Run Local (go run ./oda/tmf632/main.go)
	// connStr := "postgresql://myuser:mypass@localhost:5432/go_oda?sslmode=disable"
	// Run on Docker (docker-compose up --build -d)
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

// listIndividual retrieves a individual
func listIndividual(c *gin.Context) {
	var individuals []Individual
	query := `SELECT par.id,par.href,ind.type,par.type AS "baseType",ind.gender,ind.countryOfBirth,ind.nationality,ind.maritalStatus,ind.birthDate,ind.givenName,ind.preferredGivenName,ind.familyName,ind.legalName,ind.middleName,ind.fullName,ind.formattedName,ind.status
		FROM party par INNER JOIN individual ind ON par.id=ind.party_id;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Individual", "error": "Failed to retrieve individuals"})
		return
	}
	defer rows.Close()

	// Iterate over the result set and populate the slice
	for rows.Next() {
		var individual Individual
		if err := rows.Scan(&individual.ID, &individual.Href, &individual.Type, &individual.BaseType, &individual.Gender, &individual.CountryOfBirth, &individual.Nationality, &individual.MaritalStatus, &individual.BirthDate, &individual.GivenName, &individual.PreferredGivenName, &individual.FamilyName, &individual.LegalName, &individual.MiddleName, &individual.FullName, &individual.FormattedName, &individual.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"@type": "Individual", "error": "Failed to scan individual"})
			return
		}

		if errMsg := getExternalReference(&individual, *individual.ID); errMsg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"@type": "Individual", "error": errMsg})
			return
		}

		// Append to the individuals slice
		individuals = append(individuals, individual)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Individual", "error": "Error while fetching individuals"})
		return
	}

	c.JSON(http.StatusOK, individuals)
}

// getIndividualById retrieves a individual by ID
func getIndividualById(c *gin.Context) {
	id := c.Param("id")
	var individual Individual
	query := `SELECT par.id,par.href,ind.type,par.type AS "baseType",ind.gender,ind.countryOfBirth,ind.nationality,ind.maritalStatus,ind.birthDate,ind.givenName,ind.preferredGivenName,ind.familyName,ind.legalName,ind.middleName,ind.fullName,ind.formattedName,ind.status
		FROM party par INNER JOIN individual ind ON par.id=ind.party_id WHERE par.id = $1 LIMIT 1`
	row := db.QueryRow(query, id)
	log.Println(*row)
	if err := row.Scan(&individual.ID, &individual.Href, &individual.Type, &individual.BaseType, &individual.Gender, &individual.CountryOfBirth, &individual.Nationality, &individual.MaritalStatus, &individual.BirthDate, &individual.GivenName, &individual.PreferredGivenName, &individual.FamilyName, &individual.LegalName, &individual.MiddleName, &individual.FullName, &individual.FormattedName, &individual.Status); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"@type": "Individual", "error": "Individual not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Individual", "error": "Failed to scan individual" + err.Error()})
		return
	}

	if errMsg := getExternalReference(&individual, *individual.ID); errMsg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Individual", "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, individual)
}

// createIndividual creates a new individual
func createIndividual(c *gin.Context) {
	var newIndividual Individual
	if err := c.ShouldBindJSON(&newIndividual); err != nil {
		c.JSON(http.StatusCreated, gin.H{"@type": "Individual", "error": err.Error()})
		return
	}

	//Href
	href := "http://localhost:8081/tmf-api/party/v5/individual/" + *newIndividual.ID
	//Birthdate
	birthDate := sql.NullTime{Valid: false}
	if newIndividual.BirthDate != nil {
		// Parse the timestamp string into time.Time
		parsedTime, err := time.Parse(time.RFC3339, *newIndividual.BirthDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"@type": "Individual", "error": "Failed to parsing timestamp " + err.Error()})
			return
		}

		// Create a sql.NullTime with the parsed time
		birthDate = sql.NullTime{Time: parsedTime, Valid: true}
	}

	query := `
		WITH partyins AS (
			INSERT INTO party (id,href,type) VALUES ($1,$2,'Party') RETURNING id
		), individualins AS (
			INSERT INTO individual (gender,countryOfBirth,nationality,maritalStatus,birthDate,givenName,preferredGivenName,familyName,legalName,middleName,fullName,formattedName,status,type,party_id) VALUES ($3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,'Invididual',$1) RETURNING id
		)
		SELECT id FROM partyins;
	`
	err := db.QueryRow(query, newIndividual.ID, href, newIndividual.Gender, newIndividual.CountryOfBirth, newIndividual.Nationality, newIndividual.MaritalStatus, birthDate, newIndividual.GivenName, newIndividual.PreferredGivenName, newIndividual.FamilyName, newIndividual.LegalName, newIndividual.MiddleName, newIndividual.FullName, newIndividual.FormattedName, newIndividual.Status).Scan(&newIndividual.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Individual", "error": "Failed to insert individual " + err.Error()})
		return
	}

	if newIndividual.ExternalReferences != nil {
		for _, externalReference := range *newIndividual.ExternalReferences {
			var id int
			query = `INSERT INTO externalReference (name,externalIdentifierType,type,party_id) VALUES ($1,$2,$3,$4) RETURNING id;`
			err := db.QueryRow(query, externalReference.Name, externalReference.ExternalIdentifierType, externalReference.Type, newIndividual.ID).Scan(&id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"@type": "Individual", "error": "Failed to insert externalReference " + err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, newIndividual)
}

// updateIndividual updates a individual
func updateIndividual(c *gin.Context) {
	id := c.Param("id")
	var individual Individual
	if err := c.ShouldBindJSON(&individual); err != nil {
		c.JSON(http.StatusOK, gin.H{"@type": "Individual", "error": err.Error()})
		return
	}

	query := "UPDATE individual SET " // Base query
	var setClauses []string           // Slice to hold SET clauses
	var params []interface{}          // Slice to hold query parameters
	counter := 1                      // Counter for parameter placeholders ($1, $2, etc.)

	// Check each field and add to the query if non-empty
	if individual.Gender != nil {
		setClauses = append(setClauses, fmt.Sprintf("gender = $%d", counter))
		params = append(params, individual.Gender)
		counter++
	}
	if individual.CountryOfBirth != nil {
		setClauses = append(setClauses, fmt.Sprintf("countryOfBirth = $%d", counter))
		params = append(params, individual.CountryOfBirth)
		counter++
	}
	if individual.Nationality != nil {
		setClauses = append(setClauses, fmt.Sprintf("nationality = $%d", counter))
		params = append(params, individual.Nationality)
		counter++
	}
	if individual.MaritalStatus != nil {
		setClauses = append(setClauses, fmt.Sprintf("maritalStatus = $%d", counter))
		params = append(params, individual.MaritalStatus)
		counter++
	}
	if individual.BirthDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("birthDate = $%d", counter))
		params = append(params, individual.BirthDate)
		counter++
	}
	if individual.GivenName != nil {
		setClauses = append(setClauses, fmt.Sprintf("givenName = $%d", counter))
		params = append(params, individual.GivenName)
		counter++
	}
	if individual.PreferredGivenName != nil {
		setClauses = append(setClauses, fmt.Sprintf("preferredGivenName = $%d", counter))
		params = append(params, individual.PreferredGivenName)
		counter++
	}
	if individual.FamilyName != nil {
		setClauses = append(setClauses, fmt.Sprintf("familyName = $%d", counter))
		params = append(params, individual.FamilyName)
		counter++
	}
	if individual.LegalName != nil {
		setClauses = append(setClauses, fmt.Sprintf("legalName = $%d", counter))
		params = append(params, individual.LegalName)
		counter++
	}
	if individual.MiddleName != nil {
		setClauses = append(setClauses, fmt.Sprintf("middleName = $%d", counter))
		params = append(params, individual.MiddleName)
		counter++
	}
	if individual.FullName != nil {
		setClauses = append(setClauses, fmt.Sprintf("fullName = $%d", counter))
		params = append(params, individual.FullName)
		counter++
	}
	if individual.FormattedName != nil {
		setClauses = append(setClauses, fmt.Sprintf("formattedName = $%d", counter))
		params = append(params, individual.FormattedName)
		counter++
	}
	if individual.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", counter))
		params = append(params, individual.Status)
		counter++
	}

	// If no fields to update, return an error
	if len(setClauses) == 0 {
		c.JSON(http.StatusNotModified, gin.H{"@type": "Individual", "error": "No fields to update"})
		return
	}

	// Add the SET clauses to the query
	query += strings.Join(setClauses, ", ")

	// Add the WHERE clause
	query += fmt.Sprintf(" WHERE party_id = $%d", counter)
	params = append(params, id)

	// Execute the query
	res, err := db.Exec(query, params...)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count < 1 {
				c.JSON(http.StatusNoContent, gin.H{"@type": "Individual", "error": "Individual not found"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, individual)
}

// deleteIndividualById deletes a individual by ID
func deleteIndividualById(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM externalReference WHERE party_id = $1`
	res, err := db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count >= 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Deleted %d externalReference(s) of individual: %s. ", count, id))
			}
		}
	}
	log.Println("no externalReferences of individual:", id)

	query = `
		WITH individualdel AS (
			DELETE FROM individual WHERE party_id = $1
		)
		DELETE FROM party WHERE id = $1;
	`
	res, err = db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count == 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Deleted individual: %s. ", id))
				return
			}
		}
	}
	c.JSON(http.StatusNoContent, gin.H{"@type": "Individual", "error": "Individual not found"})
}

// listOrganization retrieves a organization
func listOrganization(c *gin.Context) {
	var organizations []Organization
	query := `SELECT par.id,par.href,org.type,par.type AS "baseType",org.isLegalEntity,org.isHeadOffice,org.organizationType,org.name,org.tradingName,org.nameType,org.status
		FROM party par INNER JOIN organization org ON par.id=org.party_id;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Organization", "error": "Failed to retrieve organizations"})
		return
	}
	defer rows.Close()

	// Iterate over the result set and populate the slice
	for rows.Next() {
		var organization Organization
		if err := rows.Scan(&organization.ID, &organization.Href, &organization.Type, &organization.BaseType, &organization.IsLegalEntity, &organization.IsHeadOffice, &organization.OrganizationType, &organization.Name, &organization.TradingName, &organization.NameType, &organization.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"@type": "Organization", "error": "Failed to scan organization"})
			return
		}

		if errMsg := getExternalReference(&organization, *organization.ID); errMsg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"@type": "Organization", "error": errMsg})
			return
		}

		// Append to the organizations slice
		organizations = append(organizations, organization)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Organization", "error": "Error while fetching organizations"})
		return
	}

	c.JSON(http.StatusOK, organizations)
}

// getOrganizationById retrieves a organization by ID
func getOrganizationById(c *gin.Context) {
	id := c.Param("id")
	var organization Organization
	query := `SELECT par.id,par.href,org.type,par.type AS "baseType",org.isLegalEntity,org.isHeadOffice,org.organizationType,org.name,org.tradingName,org.nameType,org.status
		FROM party par INNER JOIN organization org ON par.id=org.party_id WHERE par.id = $1 LIMIT 1`
	row := db.QueryRow(query, id)
	if err := row.Scan(&organization.ID, &organization.Href, &organization.Type, &organization.BaseType, &organization.IsLegalEntity, &organization.IsHeadOffice, &organization.OrganizationType, &organization.Name, &organization.TradingName, &organization.NameType, &organization.Status); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"@type": "Organization", "error": "Organization not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Organization", "error": "Failed to scan organization"})
		return
	}

	if errMsg := getExternalReference(&organization, *organization.ID); errMsg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Organization", "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, organization)
}

// createOrganization creates a new organization
func createOrganization(c *gin.Context) {
	var newOrganization Organization
	if err := c.ShouldBindJSON(&newOrganization); err != nil {
		c.JSON(http.StatusCreated, gin.H{"@type": "Organization", "error": err.Error()})
		return
	}

	//Href
	href := "http://localhost:8081/tmf-api/party/v5/organization/" + *newOrganization.ID

	query := `
		WITH partyins AS (
			INSERT INTO party (id,href,type) VALUES ($1,$2,'Party') RETURNING id
		), organizationins AS (
			INSERT INTO organization (isLegalEntity,isHeadOffice,organizationType,name,tradingName,nameType,status,type,party_id) VALUES ($3,$4,$5,$6,$7,$8,$9,'Organization',$1) RETURNING id
		)
		SELECT id FROM partyins;
	`
	err := db.QueryRow(query, newOrganization.ID, href, newOrganization.IsLegalEntity, newOrganization.IsHeadOffice, newOrganization.OrganizationType, newOrganization.Name, newOrganization.TradingName, newOrganization.NameType, newOrganization.Status).Scan(&newOrganization.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Organization", "error": "Failed to insert organization " + err.Error()})
		return
	}

	if newOrganization.ExternalReferences != nil {
		for _, externalReference := range *newOrganization.ExternalReferences {
			var id int
			query = `INSERT INTO externalReference (name,externalIdentifierType,type,party_id) VALUES ($1,$2,$3,$4) RETURNING id;`
			err := db.QueryRow(query, externalReference.Name, externalReference.ExternalIdentifierType, externalReference.Type, newOrganization.ID).Scan(&id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"@type": "Organization", "error": "Failed to insert externalReference " + err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, newOrganization)
}

// updateOrganization updates a organization
func updateOrganization(c *gin.Context) {
	id := c.Param("id")
	var organization Organization
	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusOK, gin.H{"@type": "Organization", "error": err.Error()})
		return
	}
	query := "UPDATE organization SET " // Base query
	var setClauses []string             // Slice to hold SET clauses
	var params []interface{}            // Slice to hold query parameters
	counter := 1                        // Counter for parameter placeholders ($1, $2, etc.)

	// Check each field and add to the query if non-empty
	if organization.IsLegalEntity != nil {
		setClauses = append(setClauses, fmt.Sprintf("isLegalEntity = $%d", counter))
		params = append(params, *organization.IsLegalEntity)
		counter++
	}
	if organization.IsHeadOffice != nil {
		setClauses = append(setClauses, fmt.Sprintf("isHeadOffice = $%d", counter))
		params = append(params, *organization.IsHeadOffice)
		counter++
	}
	if organization.OrganizationType != nil {
		setClauses = append(setClauses, fmt.Sprintf("organizationType = $%d", counter))
		params = append(params, organization.OrganizationType)
		counter++
	}
	if organization.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", counter))
		params = append(params, organization.Name)
		counter++
	}
	if organization.TradingName != nil {
		setClauses = append(setClauses, fmt.Sprintf("tradingName = $%d", counter))
		params = append(params, organization.TradingName)
		counter++
	}
	if organization.NameType != nil {
		setClauses = append(setClauses, fmt.Sprintf("nameType = $%d", counter))
		params = append(params, organization.NameType)
		counter++
	}
	if organization.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", counter))
		params = append(params, organization.Status)
		counter++
	}

	// If no fields to update, return an error
	if len(setClauses) == 0 {
		c.JSON(http.StatusNotModified, gin.H{"@type": "Organization", "error": "No fields to update"})
		return
	}

	// Add the SET clauses to the query
	query += strings.Join(setClauses, ", ")

	// Add the WHERE clause
	query += fmt.Sprintf(" WHERE party_id = $%d", counter)
	params = append(params, id)

	// Execute the query
	res, err := db.Exec(query, params...)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count < 1 {
				c.JSON(http.StatusNoContent, gin.H{"@type": "Organization", "error": "Organization not found"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, organization)
}

// deleteOrganizationById deletes a organization by ID
func deleteOrganizationById(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM externalReference WHERE party_id = $1`
	res, err := db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count >= 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Deleted %d externalReference(s) of organization: %s. ", count, id))
			}
		}
	}
	log.Println("no externalReferences of organization:", id)

	query = `
		WITH organizationdel AS (
			DELETE FROM organization WHERE party_id = $1
		)
		DELETE FROM party WHERE id = $1;
	`
	res, err = db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count == 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Deleted organization: %s. ", id))
				return
			}
		}
	}
	c.JSON(http.StatusNoContent, gin.H{"@type": "Organization", "error": "Organization not found"})
}

func getExternalReference(input interface{}, party_id string) string {
	var externalReferences []ExternalReference
	query := `SELECT name, externalIdentifierType, type FROM externalReference WHERE party_id = $1`
	rows, err := db.Query(query, party_id)
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
		switch v := input.(type) {
		case *Individual:
			v.ExternalReferences = &externalReferences
		case *Organization:
			v.ExternalReferences = &externalReferences
		default:
			fmt.Println("getExternalReference : Unsupported type")
		}
	}
	return ""
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/tmf-api/partyManagement/v5/individual", listIndividual)
	r.GET("/tmf-api/partyManagement/v5/individual/:id", getIndividualById)
	r.POST("/tmf-api/partyManagement/v5/individual", createIndividual)
	r.PATCH("/tmf-api/partyManagement/v5/individual/:id", updateIndividual)
	r.DELETE("/tmf-api/partyManagement/v5/individual/:id", deleteIndividualById)

	r.GET("/tmf-api/partyManagement/v5/organization", listOrganization)
	r.GET("/tmf-api/partyManagement/v5/organization/:id", getOrganizationById)
	r.POST("/tmf-api/partyManagement/v5/organization", createOrganization)
	r.PATCH("/tmf-api/partyManagement/v5/organization/:id", updateOrganization)
	r.DELETE("/tmf-api/partyManagement/v5/organization/:id", deleteOrganizationById)
	r.Run(":8081")
}
