package main

import (
	"log"
	"log/slog"
	"os"

	"io"
	"tickets/config"
	"tickets/controllers"
	db "tickets/db/sqlc"
	"tickets/publish"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/natefinch/lumberjack"
)

func initLogger() *slog.Logger {
	// Configure file logging with lumberjack
	logFile := &lumberjack.Logger{
		Filename:   "./logger/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // compress with .gz
	}

	// Create a multi-writer that writes to both file and stdout
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	// Create a single handler with the multi-writer
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
func main() {

	// Setup logger
	gin.SetMode(gin.ReleaseMode)
	initLogger()

	// Initialize DB connection
	dbConn, err := config.DBConnection()

	// Connect DB
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
		slog.Error("Failed to connect DB", "error", err)
	}
	defer dbConn.Close()

	err = publish.InitRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		slog.Error("failed to connect to rabbitmq", "error", err)
		log.Fatal("failed to connect to rabbitmq:", err)
	}

	queries := db.New(dbConn)
	tc := &controllers.TicketController{Queries: queries, DB: dbConn}
	uc := &controllers.UserController{Queries: queries, DB: dbConn}

	// Setup Gin
	r := gin.Default()

	// Ticket routes
	r.POST("/tickets", tc.CreateTicket)
	r.GET("/tickets", tc.ListTickets)
	r.GET("/tickets/:id", tc.GetTicket)
	r.PUT("/tickets/:id/status", tc.UpdateTicketStatus)
	r.POST("/users", uc.CreateUser)
	r.GET("/users", uc.ListUsers)
	r.POST("/updateuser/:id", uc.UpdateUser)

	r.Run(":8082")
}
