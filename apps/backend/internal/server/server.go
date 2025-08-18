package server

import (
	"backend/internal/actors"
	"backend/internal/broker"
	"backend/internal/database"
	"backend/internal/domains"
	"backend/internal/files"
	"backend/internal/permissions"
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

	db          database.Service
	broker      broker.Service
	actors      actors.Service
	permissions permissions.Service
	files       files.Service

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
	actorsService := actors.New(databaseService)
	filesService := files.New()
	permissionsService := permissions.New(databaseService, brokerService)

	authService := domains.NewAuthService(databaseService, brokerService)
	chatService := domains.NewChatService(actorsService, databaseService, filesService, permissionsService)
	userService := domains.NewUserService(databaseService, brokerService, filesService)
	channelService := domains.NewChannelService(databaseService, actorsService, permissionsService)
	friendService := domains.NewFriendService(databaseService, actorsService)
	roleService := domains.NewRoleService(databaseService, actorsService, permissionsService)
	serverService := domains.NewServerService(databaseService, actorsService, filesService, permissionsService)

	NewServer := &Server{
		port: port,

		db:          databaseService,
		broker:      brokerService,
		actors:      actorsService,
		files:       filesService,
		permissions: permissionsService,

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
