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

	r.Use(middlewares.NewRateLimiter(middlewares.LimiterConfig{
		MaxRequests: 100,
		Window:      30 * time.Second,
	}))

	// pprof.Register(r)

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	api := r.Group("/api")
	auth := api.Group("/auth")

	ws := handlers.NewWSHandlers(s.actors)
	auth.GET("/ws/:user_id", ws.Setup)

	chat := handlers.NewChatHandlers(s.chat)
	auth.POST("/:user_id/chat", chat.CreateMessage)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	allHealthChecks := map[string]any{
		"db":     s.db.Health(),
		"broker": s.broker.Health(),
	}

	c.JSON(http.StatusOK, allHealthChecks)
}
