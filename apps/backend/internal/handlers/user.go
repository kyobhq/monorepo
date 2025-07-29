package handlers

import (
	"backend/internal/domains"
	"backend/internal/types"
	"backend/internal/validation"
	"net/http"

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

// func (h *userHandler) GetUser(c *gin.Context) {
// 	userID := c.Param("user_id")
//
// 	user, err := h.domain.GetUserByID(c, userID)
// 	if err != nil {
// 		err.Respond(c)
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, user)
// }

func (h *userHandler) Setup(c *gin.Context) {
	setup, err := h.domain.Setup(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, setup)
}

func (h *userHandler) UpdateAccount(c *gin.Context) {
	var body types.UpdateAccountParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	userID := c.Param("user_id")

	err := h.domain.UpdateAccount(c, userID, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *userHandler) UpdateProfile(c *gin.Context) {
	var body types.UpdateProfileParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	userID := c.Param("user_id")

	err := h.domain.UpdateProfile(c, userID, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *userHandler) UpdateAvatar(c *gin.Context) {
}
