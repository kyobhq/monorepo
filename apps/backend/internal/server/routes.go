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

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	api := r.Group("/api")
	protected := api.Group("/auth")
	protected.Use(middlewares.Auth(s.broker))

	auth := handlers.NewAuthHandlers(s.auth)
	api.POST("/signin", auth.SignIn)
	api.POST("/signup", auth.SignUp)
	api.POST("/logout", auth.Logout)

	ws := handlers.NewWSHandlers(s.actors)
	protected.GET("/ws/:user_id", ws.Setup)

	chat := handlers.NewChatHandlers(s.chat)
	protected.POST("/:user_id/chat", chat.CreateMessage)
	protected.GET("/hello", s.HelloWorldHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]any)
	user, exists := c.Get("user")
	if exists {
		resp["message"] = user
	} else {
		resp["message"] = "Hello world"
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	allHealthChecks := map[string]any{
		"db":     s.db.Health(),
		"broker": s.broker.Health(),
	}

	c.JSON(http.StatusOK, allHealthChecks)
}
