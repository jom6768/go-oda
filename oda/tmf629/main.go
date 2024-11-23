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

type Customer struct {
	Type           *string          `json:"@type" binding:"required"` //"Customer"
	BaseType       *string          `json:"@baseType,omitempty"`      //"PartyRole","Producer","Consumer","Customer","BusinessPartner","Supplier"
	Href           *string          `json:"href,omitempty"`
	ID             *string          `json:"id" binding:"required"`
	Name           *string          `json:"name,omitempty"`
	Description    *string          `json:"description,omitempty"`
	Role           *string          `json:"role,omitempty"`
	Status         *string          `json:"status,omitempty"`
	StatusReason   *string          `json:"statusReason,omitempty"`
	ValidFor       *ValidFor        `json:"validFor,omitempty"`
	ContactMediums *[]ContactMedium `json:"contactMedium,omitempty"`
}

type ValidFor struct {
	StartDateTime *string `json:"startDateTime,omitempty"`
	EndDateTime   *string `json:"endDateTime,omitempty"`
}

type ContactMedium struct {
	Type        *string   `json:"@type" binding:"required"` //"PhoneContactMedium","ContactMedium","GeographicAddressContactMedium","SocialContactMedium","EmailContactMedium","FaxContactMedium"
	Preferred   *bool     `json:"preferred,omitempty"`
	ContactType *string   `json:"contactType,omitempty"`
	ValidFor    *ValidFor `json:"validFor,omitempty"`
	PhoneNumber *string   `json:"phoneNumber,omitempty"`
	City        *string   `json:"city,omitempty"`
	Country     *string   `json:"country,omitempty"`
	PostCode    *string   `json:"postCode,omitempty"`
	Street1     *string   `json:"street1,omitempty"`
}

// ////////////////////////////////////////////////
// Database Connection
// ////////////////////////////////////////////////
var db *sql.DB

func initDB() {
	var err error
	// Run Local (go run ./oda/tmf629/main.go)
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

// ////////////////////////////////////////////////
// Customer Function
// ////////////////////////////////////////////////

// listCustomer retrieves a customer
func listCustomer(c *gin.Context) {
	var customers []Customer
	query := `SELECT par.id,par.href,cus.type,par.name,par.description,par.role,par.status,par.statusReason,par.startDateTime,par.endDateTime
		FROM partyRole par INNER JOIN customer cus ON par.id=cus.partyRole_id;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Failed to retrieve customers"})
		return
	}
	defer rows.Close()

	// Iterate over the result set and populate the slice
	for rows.Next() {
		var customer Customer
		var validFor ValidFor
		if err := rows.Scan(&customer.ID, &customer.Href, &customer.Type, &customer.Name, &customer.Description, &customer.Role, &customer.Status, &customer.StatusReason, &validFor.StartDateTime, &validFor.EndDateTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Failed to scan customer"})
			return
		}
		if validFor.StartDateTime != nil || validFor.EndDateTime != nil {
			customer.ValidFor = &validFor
		}

		if errMsg := getContactMedium(&customer, *customer.ID); errMsg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": errMsg})
			return
		}

		// Append to the customers slice
		customers = append(customers, customer)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Error while fetching customers"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// getCustomerById retrieves a customer by ID
func getCustomerById(c *gin.Context) {
	id := c.Param("id")
	var customer Customer
	var validFor ValidFor
	query := `SELECT par.id,par.href,cus.type,par.name,par.description,par.role,par.status,par.statusReason,par.startDateTime,par.endDateTime
		FROM partyRole par INNER JOIN customer cus ON par.id=cus.partyRole_id WHERE par.id = $1 LIMIT 1;`
	row := db.QueryRow(query, id)
	if err := row.Scan(&customer.ID, &customer.Href, &customer.Type, &customer.Name, &customer.Description, &customer.Role, &customer.Status, &customer.StatusReason, &validFor.StartDateTime, &validFor.EndDateTime); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"@type": "Customer", "error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Failed to scan customer" + err.Error()})
		return
	}
	if validFor.StartDateTime != nil || validFor.EndDateTime != nil {
		customer.ValidFor = &validFor
	}

	if errMsg := getContactMedium(&customer, *customer.ID); errMsg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": errMsg})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// createCustomer creates a new customer
func createCustomer(c *gin.Context) {
	var newCustomer Customer
	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusCreated, gin.H{"@type": "Customer", "error": err.Error()})
		return
	}

	//Href
	href := "http://localhost:8629/tmf-api/customerManagement/v5/customer/" + *newCustomer.ID
	//ValidFor
	startDateTime := sql.NullTime{Valid: false}
	endDateTime := sql.NullTime{Valid: false}
	if newCustomer.ValidFor != nil {
		//ValidFor.StartDateTime
		if newCustomer.ValidFor.StartDateTime != nil {
			// Parse the timestamp string into time.Time
			parsedTime, err := time.Parse(time.RFC3339, *newCustomer.ValidFor.StartDateTime)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Failed to parsing timestamp " + err.Error()})
				return
			}

			// Create a sql.NullTime with the parsed time
			startDateTime = sql.NullTime{Time: parsedTime, Valid: true}
		}
		//ValidFor.EndDateTime
		if newCustomer.ValidFor.EndDateTime != nil {
			// Parse the timestamp string into time.Time
			parsedTime, err := time.Parse(time.RFC3339, *newCustomer.ValidFor.EndDateTime)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Failed to parsing timestamp " + err.Error()})
				return
			}

			// Create a sql.NullTime with the parsed time
			endDateTime = sql.NullTime{Time: parsedTime, Valid: true}
		}
	}

	query := `
		WITH partyins AS (
			INSERT INTO partyRole (id,href,name,description,role,status,statusReason,startDateTime,endDateTime,type) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,'PartyRole') RETURNING id
		), customerins AS (
			INSERT INTO customer (type,partyRole_id) VALUES ('Invididual',$1) RETURNING id
		)
		SELECT id FROM partyins;
	`
	err := db.QueryRow(query, newCustomer.ID, href, newCustomer.Name, newCustomer.Description, newCustomer.Role, newCustomer.Status, newCustomer.StatusReason, startDateTime, endDateTime).Scan(&newCustomer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Failed to insert customer " + err.Error()})
		return
	}

	if newCustomer.ContactMediums != nil {
		for _, contactMedium := range *newCustomer.ContactMediums {
			//ValidFor
			startDateTime := sql.NullTime{Valid: false}
			endDateTime := sql.NullTime{Valid: false}
			if contactMedium.ValidFor != nil {
				//ValidFor.StartDateTime
				if contactMedium.ValidFor.StartDateTime != nil {
					// Parse the timestamp string into time.Time
					parsedTime, err := time.Parse(time.RFC3339, *contactMedium.ValidFor.StartDateTime)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Failed to parsing timestamp " + err.Error()})
						return
					}

					// Create a sql.NullTime with the parsed time
					startDateTime = sql.NullTime{Time: parsedTime, Valid: true}
				}
				//ValidFor.EndDateTime
				if contactMedium.ValidFor.EndDateTime != nil {
					// Parse the timestamp string into time.Time
					parsedTime, err := time.Parse(time.RFC3339, *contactMedium.ValidFor.EndDateTime)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Failed to parsing timestamp " + err.Error()})
						return
					}

					// Create a sql.NullTime with the parsed time
					endDateTime = sql.NullTime{Time: parsedTime, Valid: true}
				}
			}
			var id int
			query = `INSERT INTO contactMedium (preferred,contactType,phoneNumber,city,country,postCode,street1,startDateTime,endDateTime,type,partyRole_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id;`
			err := db.QueryRow(query, contactMedium.Preferred, contactMedium.ContactType, contactMedium.PhoneNumber, contactMedium.City, contactMedium.Country, contactMedium.PostCode, contactMedium.Street1, startDateTime, endDateTime, contactMedium.Type, newCustomer.ID).Scan(&id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"@type": "Customer", "error": "Failed to insert engagedParty " + err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, newCustomer)
}

// updateCustomer updates a customer
func updateCustomer(c *gin.Context) {
	id := c.Param("id")
	var customer Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusOK, gin.H{"@type": "Customer", "error": err.Error()})
		return
	}

	query := "UPDATE partyRole SET " // Base query
	var setClauses []string          // Slice to hold SET clauses
	var params []interface{}         // Slice to hold query parameters
	counter := 1                     // Counter for parameter placeholders ($1, $2, etc.)

	// Check each field and add to the query if non-empty
	if customer.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", counter))
		params = append(params, customer.Name)
		counter++
	}
	if customer.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", counter))
		params = append(params, *customer.Description)
		counter++
	}
	if customer.Role != nil {
		setClauses = append(setClauses, fmt.Sprintf("role = $%d", counter))
		params = append(params, *customer.Role)
		counter++
	}
	if customer.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", counter))
		params = append(params, *customer.Status)
		counter++
	}
	if customer.StatusReason != nil {
		setClauses = append(setClauses, fmt.Sprintf("statusReason = $%d", counter))
		params = append(params, *customer.StatusReason)
		counter++
	}
	if customer.ValidFor != nil {
		if customer.ValidFor.StartDateTime != nil {
			setClauses = append(setClauses, fmt.Sprintf("startDateTime = $%d", counter))
			params = append(params, *customer.ValidFor.StartDateTime)
			counter++
		}
		if customer.ValidFor.EndDateTime != nil {
			setClauses = append(setClauses, fmt.Sprintf("endDateTime = $%d", counter))
			params = append(params, *customer.ValidFor.EndDateTime)
			counter++
		}
	}

	// If no fields to update, return an error
	if len(setClauses) == 0 {
		c.JSON(http.StatusNotModified, gin.H{"@type": "Customer", "error": "No fields to update"})
		return
	}

	// Add the SET clauses to the query
	query += strings.Join(setClauses, ", ")

	// Add the WHERE clause
	query += fmt.Sprintf(" WHERE id = $%d", counter)
	params = append(params, id)

	// Execute the query
	res, err := db.Exec(query, params...)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count < 1 {
				c.JSON(http.StatusNoContent, gin.H{"@type": "Customer", "error": "Customer not found"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, customer)
}

// deleteCustomerById deletes a customer by ID
func deleteCustomerById(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM engagedParty WHERE partyRole_id = $1`
	res, err := db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count >= 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Deleted %d engagedParty(s) of customer: %s. ", count, id))
			}
		}
	}

	query = `DELETE FROM account WHERE partyRole_id = $1`
	res, err = db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count >= 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Deleted %d account(s) of customer: %s. ", count, id))
			}
		}
	}

	query = `DELETE FROM paymentMethod WHERE partyRole_id = $1`
	res, err = db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count >= 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Deleted %d paymentMethod(s) of customer: %s. ", count, id))
			}
		}
	}

	query = `DELETE FROM contactMedium WHERE partyRole_id = $1`
	res, err = db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count >= 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Deleted %d contactMedium(s) of customer: %s. ", count, id))
			}
		}
	}

	query = `
		WITH customerdel AS (
			DELETE FROM customer WHERE partyRole_id = $1
		)
		DELETE FROM partyRole WHERE id = $1;
	`
	res, err = db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count == 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("Deleted customer: %s. ", id))
				return
			}
		}
	}
	c.JSON(http.StatusNoContent, gin.H{"@type": "Customer", "error": "Customer not found"})
}

// ////////////////////////////////////////////////
// Contact Medium Function
// ////////////////////////////////////////////////

func getContactMedium(input interface{}, partyRole_id string) string {
	var contactMediums []ContactMedium
	query := `SELECT preferred,contactType,phoneNumber,city,country,postCode,street1,startDateTime,endDateTime,type FROM contactMedium WHERE partyRole_id = $1;`
	rows, err := db.Query(query, partyRole_id)
	if err != nil {
		return "Failed to retrieve contactMediums"
	}
	defer rows.Close()

	// Iterate over the result set and populate the slice
	for rows.Next() {
		var contactMedium ContactMedium
		var validFor ValidFor
		if err := rows.Scan(&contactMedium.Preferred, &contactMedium.ContactType, &contactMedium.PhoneNumber, &contactMedium.City, &contactMedium.Country, &contactMedium.PostCode, &contactMedium.Street1, &validFor.StartDateTime, &validFor.EndDateTime, &contactMedium.Type); err != nil {
			return "Failed to scan contactMedium"
		}
		if validFor.StartDateTime != nil || validFor.EndDateTime != nil {
			contactMedium.ValidFor = &validFor
		}

		// Append to the contactMediums slice
		contactMediums = append(contactMediums, contactMedium)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return "Error while fetching contactMediums"
	}

	if contactMediums != nil {
		switch v := input.(type) {
		case *Customer:
			v.ContactMediums = &contactMediums
		default:
			fmt.Println("getContactMedium : Unsupported type")
		}
	}
	return ""
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/tmf-api/customerManagement/v5/customer", listCustomer)
	r.GET("/tmf-api/customerManagement/v5/customer/:id", getCustomerById)
	r.POST("/tmf-api/customerManagement/v5/customer", createCustomer)
	r.PATCH("/tmf-api/customerManagement/v5/customer/:id", updateCustomer)
	r.DELETE("/tmf-api/customerManagement/v5/customer/:id", deleteCustomerById)
	r.Run(":8629")
}
