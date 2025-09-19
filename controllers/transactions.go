package controllers

import (
	"database/sql"
	// "encoding/json"
	"net/http"
	"strconv"

	"log/slog"
	db "tickets/db/sqlc"

	// "tickets/publish"

	"github.com/gin-gonic/gin"
)

type TransactionsController struct {
	Queries *db.Queries
	DB      *sql.DB
}

func (ct *TransactionsController) ListTransactions(c *gin.Context) {
	// Example: parse limit and offset from query params, or set defaults
	limit := int32(10)
	offset := int32(0)
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "10")); err == nil {
		limit = int32(l)
	}
	if o, err := strconv.Atoi(c.DefaultQuery("offset", "0")); err == nil {
		offset = int32(o)
	}

	params := db.ListTransactionsParams{
		Limit:  limit,
		Offset: offset,
	}

	transactions, err := ct.Queries.ListTransactions(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list transactions"})
		slog.Error("Failed to list transactions", "error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
	slog.Info("transactions listed successfully", "count", len(transactions))
}
func (ct *TransactionsController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		slog.Error("Invalid transaction ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	transaction, err := ct.Queries.GetTransanctionByID(c.Request.Context(), int32(id))
	if err != nil {
		// slog.Error("Transaction not found", "error", err)
		slog.Error(err.Error(), "error", "not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
	slog.Info("Fetched transaction successfully", "transaction_id", transaction.ID)
}

type CreateTransactionRequest struct {
	TransactionID string  `json:"transaction_id"`
	UserID        int64   `json:"user_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        int     `json:"status"`
	PaymentMethod string  `json:"payment_method"`
}

func (ct *TransactionsController) CreateTransactions(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params := db.CreateTransactionParams{
		TransactionID: req.TransactionID,
		UserID:        int32(req.UserID),
		Amount:        strconv.FormatFloat(req.Amount, 'f', -1, 64), // Convert float64 to string for Amount
		Currency:      req.Currency,
		Status:        int16(req.Status),
		PaymentMethod: sql.NullString{
			String: req.PaymentMethod,
			Valid:  true,
		},
	}

	transaction, err := ct.Queries.CreateTransaction(c.Request.Context(), params)
	if err != nil {
		slog.Error("Failed to create transaction", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create transaction"})
		return
	}

	id, _ := transaction.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transaction created",
		"id":          id,
		"transaction": transaction,
	})
	slog.Info("Transaction created successfully", "id", id)
}
