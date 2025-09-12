package controllers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	db "tickets/db/sqlc"
	"tickets/publish"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	Queries *db.Queries
	DB      *sql.DB
}

// Create Ticket
type CreateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedBy   int64  `json:"created_by"`
	Priority    string `json:"priority"`
	Status      int16  `json:"status"` // should be int16, not string
}

func (t *TicketController) CreateTicket(c *gin.Context) {
	var req CreateTicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("Invalid request payload", "error", err)
		return
	}

	ticket, err := t.Queries.CreateTicket(c, db.CreateTicketParams{
		Title:       req.Title,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
		Priority:    req.Priority,
		Status:      req.Status, // now matches int16
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ticket"})
		slog.Error("Failed to create ticket", "error", err)
		return
	}
	// publish to RabbitMQ
	event := map[string]interface{}{
		"type":    "ticket.created",
		"payload": ticket,
	}
	body, _ := json.Marshal(event)
	if err := publish.Publish("ticket_events", body); err != nil {
		slog.Error("Failed to publish ticket event", "error", err)
	}
	c.JSON(http.StatusOK, gin.H{"ticket_id": ticket})
	slog.Info("Ticket created successfully", "ticket_id", ticket)
}

// List Tickets
func (tc *TicketController) ListTickets(c *gin.Context) {
	tickets, err := tc.Queries.ListTickets(c.Request.Context(), db.ListTicketsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		slog.Error("Failed to fetch tickets", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
	slog.Info("Fetched tickets successfully", "count", len(tickets))
}

// Get Ticket
func (tc *TicketController) GetTicket(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	ticket, err := tc.Queries.GetTicket(c.Request.Context(), id)
	if err != nil {
		slog.Error("Ticket not found", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
	slog.Info("Fetched ticket successfully", "ticket_id", ticket.ID)
}

// Update Ticket Status
func (tc *TicketController) UpdateTicketStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("Invalid request payload", "error", err)
		return
	}

	// convert string -> int16
	statusInt, err := strconv.ParseInt(req.Status, 10, 16)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	err = tc.Queries.UpdateTicketStatus(c.Request.Context(), db.UpdateTicketStatusParams{
		Status: int16(statusInt), // now matches sqlc expectation
		ID:     id,
	})
	if err != nil {
		slog.Error("Failed to update ticket status", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket status updated"})
}
