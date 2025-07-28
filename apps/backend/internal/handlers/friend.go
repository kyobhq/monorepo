package handlers

import (
	"backend/internal/domains"

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
}

func (h *friendHandler) AcceptRequest(c *gin.Context) {
}

func (h *friendHandler) RemoveFriend(c *gin.Context) {
}
