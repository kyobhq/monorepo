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
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Use(middlewares.RateLimiter(middlewares.LimiterConfig{
		MaxRequests: 100,
		Window:      30 * time.Second,
	}))

	// pprof.Register(r)

	r.GET("/health", s.healthHandler)

	api := r.Group("/api")
	protected := api.Group("/protected")
	protected.Use(middlewares.Auth(s.broker))

	auth := handlers.NewAuthHandlers(s.authSvc)
	api.POST("/signin", auth.SignIn)
	api.POST("/signup", auth.SignUp)
	protected.POST("/logout", auth.Logout)

	ws := handlers.NewWSHandlers(s.actors)
	protected.GET("/ws/:user_id", ws.Setup)

	user := handlers.NewUserHandlers(s.userSvc)
	protected.GET("/users/:user_id", user.GetUserProfile)
	protected.GET("/users/setup", user.Setup)
	protected.PATCH("/users/email", user.UpdateEmail)
	protected.PATCH("/users/password", user.UpdatePassword)
	protected.PATCH("/users/profile", user.UpdateProfile)
	protected.PATCH("/users/avatar", user.UpdateAvatar)
	protected.POST("/users/emojis", user.UploadEmojis)
	protected.PATCH("/users/emojis/:emoji_id", user.UpdateEmojis)
	protected.DELETE("/users/emojis/:emoji_id", user.DeleteEmoji)
	protected.DELETE("/users", user.DeleteAccount)
	protected.POST("/users/sync", user.Sync)

	friend := handlers.NewFriendHandlers(s.friendSvc)
	protected.POST("/friends", friend.SendRequest)
	protected.PATCH("/friends", friend.AcceptRequest)
	protected.DELETE("/friends", friend.RemoveFriend)

	server := handlers.NewServerHandlers(s.serverSvc)
	protected.POST("/servers", server.CreateServer)
	protected.GET("/servers/:server_id", server.GetInformations)
	protected.GET("/servers/:server_id/members", server.GetMembers)
	protected.GET("/servers/:server_id/bans", server.GetBannedMembers)
	protected.GET("/servers/:server_id/search", server.SearchMembers)
	protected.POST("/servers/join", server.JoinServer)
	protected.POST("/servers/:server_id/leave", server.LeaveServer)
	protected.POST("/servers/:server_id/invite", server.CreateInvite)
	protected.POST("/servers/:server_id/ban", server.BanUser)
	protected.POST("/servers/:server_id/unban/:user_id", server.UnbanUser)
	protected.POST("/servers/:server_id/kick", server.KickUser)
	protected.DELETE("/servers/invite/:invite_id", server.DeleteInvite)
	protected.PATCH("/servers/:server_id/profile", server.UpdateProfile)
	protected.PATCH("/servers/:server_id/avatar", server.UpdateAvatar)
	protected.DELETE("/servers/:server_id", server.DeleteServer)

	channel := handlers.NewChannelHandlers(s.channelSvc)
	protected.POST("/channels", channel.CreateChannel)
	protected.POST("/channels/category", channel.CreateCategory)
	protected.POST("/channels/pin/:channel_id", channel.PinChannel)
	protected.PATCH("/channels/category/:category_id", channel.EditCategory)
	protected.PATCH("/channels/:channel_id", channel.EditChannel)
	protected.DELETE("/channels/:channel_id", channel.DeleteChannel)
	protected.DELETE("/channels/category/:category_id", channel.DeleteCategory)

	chat := handlers.NewChatHandlers(s.chatSvc)
	protected.GET("/messages/:server_id/:channel_id", chat.GetMessages)
	protected.POST("/messages", chat.CreateMessage)
	protected.PATCH("/messages/:message_id", chat.EditMessage)
	protected.DELETE("/messages/:message_id", chat.DeleteMessage)

	role := handlers.NewRoleHandlers(s.roleSvc)
	protected.GET("/roles/:server_id", role.GetRoles)
	protected.GET("/roles/members/:role_id", role.GetRoleMembers)
	protected.POST("/roles", role.CreateOrEditRole)
	protected.PATCH("/roles/add_member", role.AddRoleMember)
	protected.PATCH("/roles/remove_member", role.RemoveRoleMember)
	protected.PATCH("/roles/move", role.MoveRole)
	protected.DELETE("/roles", role.DeleteRole)

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	allHealthChecks := map[string]any{
		"db":     s.db.Health(),
		"broker": s.broker.Health(),
	}

	c.JSON(http.StatusOK, allHealthChecks)
}
