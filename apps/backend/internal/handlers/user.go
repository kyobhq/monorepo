package handlers

import (
	"backend/internal/domains"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	domain domains.UserService
}

func NewUserHandlers(userService domains.UserService) *userHandler {
	return &userHandler{
		domain: userService,
	}
}

func (h *userHandler) GetUser(c *gin.Context) {
}

func (h *userHandler) UpdateAccount(c *gin.Context) {
}

func (h *userHandler) UpdateProfile(c *gin.Context) {
}
