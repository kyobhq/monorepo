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

	authSvc    domains.AuthService
	chatSvc    domains.ChatService
	userSvc    domains.UserService
	channelSvc domains.ChannelService
	roleSvc    domains.RoleService
	friendSvc  domains.FriendService
	serverSvc  domains.ServerService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	validation.New()
	databaseService := database.New()
	brokerService := broker.New()
	actorsService := actors.New()

	authService := domains.NewAuthService(databaseService, brokerService)
	chatService := domains.NewChatService(actorsService, databaseService)
	userService := domains.NewUserService(databaseService)
	channelService := domains.NewChannelService(databaseService)
	friendService := domains.NewFriendService(databaseService)
	roleService := domains.NewRoleService(databaseService)
	serverService := domains.NewServerService(databaseService)

	NewServer := &Server{
		port: port,

		db:     databaseService,
		broker: brokerService,
		actors: actorsService,

		authSvc:    authService,
		chatSvc:    chatService,
		userSvc:    userService,
		channelSvc: channelService,
		roleSvc:    roleService,
		friendSvc:  friendService,
		serverSvc:  serverService,
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
