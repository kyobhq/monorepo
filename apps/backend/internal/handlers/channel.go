package handlers

import (
	"backend/internal/domains"

	"github.com/gin-gonic/gin"
)

type channelHandler struct {
	domain domains.ChannelService
}

func NewChannelHandlers(channelService domains.ChannelService) *channelHandler {
	return &channelHandler{
		domain: channelService,
	}
}

func (h *channelHandler) CreateChannel(c *gin.Context) {
}

func (h *channelHandler) EditChannel(c *gin.Context) {
}

func (h *channelHandler) DeleteChannel(c *gin.Context) {
}
