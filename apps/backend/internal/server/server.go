package server

import (
	"backend/internal/actors"
	"backend/internal/broker"
	"backend/internal/database"
	"backend/internal/domains"
	"backend/internal/validation"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int

	db     database.Service
	broker broker.Service
	actors actors.Service

	chat domains.ChatService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	validation.New()
	databaseService := database.New()
	brokerService := broker.New()
	actorsService := actors.New()

	chatService := domains.NewChatService(actorsService, databaseService)

	NewServer := &Server{
		port: port,

		db:     databaseService,
		broker: brokerService,
		actors: actorsService,

		chat: chatService,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
