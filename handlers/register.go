package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	db "tickets/db/sqlc"
)

type registerReq struct {
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "full_name, phone and password are required"})
		slog.Error("invalid register request", "error", err)
		return
	}

	// ✅ Check if profile already exists
	_, err := h.queries.GetProfileByPhone(c, req.Phone)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		slog.Error("failed to check existing profile", "phone", req.Phone, "error", err)
		return
	}
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "profile already exists"})
		slog.Warn("profile already exists", "phone", req.Phone)
		return
	}

	// ✅ Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		slog.Error("failed to hash password", "error", err)
		return
	}

	// ✅ Insert into DB
	profile, err := h.queries.CreateProfile(c, db.CreateProfileParams{
		FullName:     sql.NullString{String: req.FullName, Valid: req.FullName != ""},
		Phone:        req.Phone,
		PasswordHash: string(hashed),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
		slog.Error("failed to create profile", "phone", req.Phone, "error", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "profile created successfully",
		"profile": profile,
	})
}
