package handlers

import (
	"backend/internal/domains"
	"backend/internal/types"
	"backend/internal/validation"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	domain domains.AuthService
}

func NewAuthHandlers(authService domains.AuthService) *authHandler {
	return &authHandler{
		domain: authService,
	}
}

func (h *authHandler) SignIn(c *gin.Context) {
	var body types.SignInParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	token, derr := h.domain.SignIn(c, &body)
	if derr != nil {
		derr.Respond(c)
		return
	}

	c.SetCookie("token", *token, int(time.Now().Add(30*(24*time.Hour)).Unix()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *authHandler) SignUp(c *gin.Context) {
	var body types.SignUpParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	token, derr := h.domain.SignUp(c, &body)
	if derr != nil {
		derr.Respond(c)
		return
	}

	c.SetCookie("token", *token, int(time.Now().Add(30*(24*time.Hour)).Unix()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *authHandler) Check(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *authHandler) Logout(c *gin.Context) {
	if derr := h.domain.Logout(c); derr != nil {
		derr.Respond(c)
		return
	}

	c.SetCookie("token", "", int(time.Now().Add(-30*(24*time.Hour)).Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
