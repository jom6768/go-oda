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
	ID     string `json:"customer_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"email"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

var db *sql.DB

func initDB() {
	var err error
	// Run Local
	// connStr := "postgresql://myuser:mypass@localhost:5432/go_oda?sslmode=disable"
	// Run on Docker
	connStr := "postgresql://myuser:mypass@host.docker.internal:5432/go_oda?sslmode=disable"
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

	query := `INSERT INTO customer (customer_id, name, email, phone, status) VALUES ($1, $2, $3, $4, $5) RETURNING customer_id`
	err := db.QueryRow(query, newCustomer.ID, newCustomer.Name, newCustomer.Email, newCustomer.Phone, "Active").Scan(&newCustomer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert customer" + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newCustomer)
}

// getCustomer retrieves a customer by ID
func getCustomerById(c *gin.Context) {
	id := c.Param("id")
	var customer Customer

	query := `SELECT customer_id, name, email, phone, status FROM customer WHERE customer_id = $1 LIMIT 1`
	row := db.QueryRow(query, id)
	if err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Status); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// deleteCustomer deletes a customer by ID
func deleteCustomerById(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM customer WHERE customer_id = $1`
	res, err := db.Exec(query, id)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			if count == 1 {
				c.JSON(http.StatusOK, fmt.Sprintf("customer_id = %s completed.", id))
				return
			}
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/tmf632/customer", createCustomer)
	r.GET("/tmf632/customer/:id", getCustomerById)
	r.DELETE("/tmf632/customer/:id", deleteCustomerById)
	r.Run(":8081")
}
