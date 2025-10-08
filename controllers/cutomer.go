package controllers

import (
	"database/sql"
	"log/slog"
	"net/http"
	db "tickets/db/sqlc"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	Queries *db.Queries
	DB      *sql.DB
}

func NewCustomerController(queries *db.Queries, db *sql.DB) *CustomerController {
	return &CustomerController{
		Queries: queries,
		DB:      db,
	}
}

func (controller *CustomerController) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Invalid request payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	arg := db.CreateCustomerParams{
		FullName:    req.FullName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}
	customer, err := controller.Queries.CreateCustomer(c, arg)
	if err != nil {
		slog.Error("Failed to create customer", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

type CreateCustomerRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,min=6"`
}

func (controller *CustomerController) GetCustomers(c *gin.Context) {
	// Parse pagination parameters
	limit := 10
	offset := 0

	params := db.GetCustomersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	customers, err := controller.Queries.GetCustomers(c, params)
	if err != nil {
		slog.Error("Failed to fetch customers", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, customers)
	slog.Info("Fetched customers successfully", "count", len(customers))
}

// const {tatus, data, send, open, close } = useWebSocket('ws://websocketurl')
