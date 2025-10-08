package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	db "tickets/db/sqlc"
	"tickets/sms"
	"tickets/utils"
)

type AuthHandler struct {
	queries     *db.Queries
	smsProvider sms.SMSProvider
	otpTTL      time.Duration
	maxAttempts int
}

func NewAuthHandler(q *db.Queries, p sms.SMSProvider) *AuthHandler {
	ttlMin := 5
	if v := os.Getenv("OTP_TTL_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			ttlMin = n
		}
	}
	maxA := 5
	if v := os.Getenv("OTP_MAX_ATTEMPTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxA = n
		}
	}
	return &AuthHandler{
		queries:     q,
		smsProvider: p,
		otpTTL:      time.Duration(ttlMin) * time.Minute,
		maxAttempts: maxA,
	}
}

// helper generate 6-digit secure otp
func generateOTP() (string, error) {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	num := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	code := num % 1000000
	return fmt.Sprintf("%06d", code), nil
}

type loginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and password required"})
		slog.Error("invalid login request", "error", err)
		return
	}

	profile, err := h.queries.GetProfileByPhone(c, req.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			slog.Error("invalid login credentials", "phone", req.Phone)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		slog.Error("failed to get profile by phone", "phone", req.Phone, "error", err)
		return
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(profile.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// delete expired OTPs (cleanup)
	if err := h.queries.DeleteExpiredOTPs(c); err != nil {
		log.Println("cleanup expired otp:", err)
	}

	otp, err := generateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate otp"})
		return
	}

	expires := time.Now().Add(h.otpTTL)

	// Save OTP linked to profile
	if _, err := h.queries.CreateOTP(c, db.CreateOTPParams{
		ProfileID: profile.ID,
		OtpCode:   otp,
		ExpiresAt: expires,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save otp"})
		return
	}

	// send via SMS (async)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		msg := fmt.Sprintf("Your login OTP is %s. It expires in %d minutes.", otp, int(h.otpTTL.Minutes()))
		if err := h.smsProvider.SendSMS(ctx, profile.Phone, msg); err != nil {
			log.Println("failed to send otp sms:", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent"})
}

type verifyReq struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req verifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and otp required"})
		return
	}

	profile, err := h.queries.GetProfileByPhone(c, req.Phone)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid otp or phone"})
		return
	}

	otpRec, err := h.queries.GetLatestOTPByProfileID(c, profile.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid otp"})
		return
	}

	// check attempts
	if otpRec.Attempts >= int32(h.maxAttempts) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "max OTP attempts exceeded"})
		return
	}

	// check expiry
	if time.Now().After(otpRec.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "otp expired"})
		return
	}

	// compare
	if otpRec.OtpCode != req.OTP {
		if err := h.queries.IncrementOTPAttempts(c, otpRec.ID); err != nil {
			log.Println("failed increment attempts:", err)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid otp"})
		return
	}

	// mark verified
	if err := h.queries.MarkOTPVerified(c, otpRec.ID); err != nil {
		log.Println("failed mark verified:", err)
	}

	// issue JWT for profile
	token, err := utils.GenerateJWT(int64(otpRec.ProfileID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
}
