package handlers

import (
	"backend/internal/domains"
	"backend/internal/types"
	"backend/internal/validation"
	"net/http"

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

func (h *channelHandler) CreateCategory(c *gin.Context) {
	var body types.CreateCategoryParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	category, err := h.domain.CreateCategory(c, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *channelHandler) DeleteCategory(c *gin.Context) {
	var body types.DeleteCategoryParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.DeleteCategory(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *channelHandler) CreateChannel(c *gin.Context) {
	var body types.CreateChannelParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	channel, err := h.domain.CreateChannel(c, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, channel)
}

func (h *channelHandler) PinChannel(c *gin.Context) {
	var body types.PinChannelParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.PinChannel(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *channelHandler) EditChannel(c *gin.Context) {
	var body types.EditChannelParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.EditChannel(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *channelHandler) DeleteChannel(c *gin.Context) {
	var body types.DeleteChannelParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.DeleteChannel(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
