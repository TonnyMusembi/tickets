package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"log/slog"
	db "tickets/db/sqlc"
	"tickets/publish"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Queries *db.Queries
	DB      *sql.DB
}

// Create User
type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	// Password string `json:"password"`
}

func (u *UserController) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("Invalid request payload", "error", err)
		return
	}

	user, err := u.Queries.CreateUser(c, db.CreateUserParams{
		FullName: req.FullName,
		Email:    req.Email,
		// Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		slog.Error("Failed to create user", "error", err)
		return
	}

	payload := map[string]interface{}{
		"type": "user.created",
		"payload": map[string]string{
			"email":     req.Email,
			"full_name": req.FullName,
		},
	}
	if data, err := json.Marshal(payload); err == nil {
		publish.Publish("user_events", data)
	} else {
		slog.Error("Failed to marshal user created event", "error", err)
	}

	c.JSON(http.StatusOK, gin.H{"user": user})

	userID, err := user.LastInsertId()
	if err != nil {
		slog.Error("Failed to get last insert ID", "error", err)
	} else {
		slog.Info("User created successfully", "user_id", userID)
	}

}

// List Users
func (uc *UserController) ListUsers(c *gin.Context) {
	users, err := uc.Queries.ListUsers(c.Request.Context(), db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		slog.Error("Failed to list users", "error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
	slog.Info("Users listed successfully", "count", len(users))
}

func (u *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`

		// Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("Invalid request payload", "error", err)
		return
	}
	err = u.Queries.UpdateUser(c, db.UpdateUserParams{
		ID:       id,
		FullName: req.FullName,
		Email:    req.Email,
		// Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		slog.Error("Failed to update user", "error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
	slog.Info("User updated successfully", "user_id", id)
}

func UpdateUser(c *gin.Context) {

	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Implement the update logic here
	

}
