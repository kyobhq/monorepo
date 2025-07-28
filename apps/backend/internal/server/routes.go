package server

import (
	"backend/internal/handlers"
	"backend/internal/middlewares"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.Use(middlewares.RateLimiter(middlewares.LimiterConfig{
		MaxRequests: 100,
		Window:      30 * time.Second,
	}))

	// pprof.Register(r)

	r.GET("/health", s.healthHandler)

	api := r.Group("/api")
	protected := api.Group("/auth")
	protected.Use(middlewares.Auth(s.broker))

	auth := handlers.NewAuthHandlers(s.authSvc)
	api.POST("/signin", auth.SignIn)
	api.POST("/signup", auth.SignUp)
	api.POST("/logout", auth.Logout)

	ws := handlers.NewWSHandlers(s.actors)
	protected.GET("/ws/:user_id", ws.Setup)

	user := handlers.NewUserHandlers(s.userSvc)
	protected.GET("/users/:user_id", user.GetUser)
	protected.PATCH("/users/:user_id/account", user.UpdateAccount)
	protected.PATCH("/users/:user_id/profile", user.UpdateProfile)

	friend := handlers.NewFriendHandlers(s.friendSvc)
	protected.POST("/friends/:receiver_id", friend.SendRequest)
	protected.PATCH("/friends/:sender_id", friend.AcceptRequest)
	protected.DELETE("/friends/:friendship_id", friend.RemoveFriend)

	server := handlers.NewServerHandlers(s.serverSvc)
	protected.POST("/servers", server.CreateServer)
	protected.POST("/servers/:server_id/join", server.JoinServer)
	protected.POST("/servers/:server_id/leave", server.LeaveServer)
	protected.POST("/servers/:server_id/invite", server.CreateInvite)
	protected.DELETE("/servers/:server_id/invite", server.DeleteInvite)
	protected.PATCH("/servers/:server_id/profile", server.EditProfile)
	protected.DELETE("/servers/:server_id", server.DeleteServer)

	channel := handlers.NewChannelHandlers(s.channelSvc)
	protected.POST("/channel", channel.CreateChannel)
	protected.PATCH("/channel/:channel_id", channel.EditChannel)
	protected.DELETE("/channel/:channel_id", channel.DeleteChannel)

	chat := handlers.NewChatHandlers(s.chatSvc)
	protected.GET("/messages/:server_id", chat.GetMessages)
	protected.POST("/messages", chat.CreateMessage)
	protected.PATCH("/messages/:message_id", chat.EditMessage)
	protected.DELETE("/messages/:message_id", chat.DeleteMessage)

	role := handlers.NewRoleHandlers(s.roleSvc)
	protected.GET("/roles/:server_id", role.GetRoles)
	protected.POST("/roles/:server_id", role.CreateRole)
	protected.PATCH("/roles/:role_id", role.EditRole)
	protected.PATCH("/roles/:role_id/add_member", role.AddRoleMember)
	protected.PATCH("/roles/:role_id/remove_member", role.RemoveRoleMember)
	protected.PATCH("/roles/:role_id/move", role.MoveRole)
	protected.DELETE("/roles/:role_id", role.DeleteRole)

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	allHealthChecks := map[string]any{
		"db":     s.db.Health(),
		"broker": s.broker.Health(),
	}

	c.JSON(http.StatusOK, allHealthChecks)
}
