package handlers

import (
	"backend/internal/domains"
	"backend/internal/types"
	"backend/internal/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type friendHandler struct {
	domain domains.FriendService
}

func NewFriendHandlers(friendService domains.FriendService) *friendHandler {
	return &friendHandler{
		domain: friendService,
	}
}

func (h *friendHandler) SendRequest(c *gin.Context) {
	var body types.SendRequestParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	friendship, err := h.domain.SendRequest(c, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, friendship)
}

func (h *friendHandler) AcceptRequest(c *gin.Context) {
	var body types.AcceptRequestParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	err := h.domain.AcceptRequest(c, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *friendHandler) RemoveFriend(c *gin.Context) {
	var body types.RemoveFriendParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.RemoveFriend(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
