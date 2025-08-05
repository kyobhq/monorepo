package handlers

import (
	"backend/internal/domains"
	"backend/internal/types"
	"backend/internal/validation"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type chatHandler struct {
	domain domains.ChatService
}

func NewChatHandlers(chatService domains.ChatService) *chatHandler {
	return &chatHandler{
		domain: chatService,
	}
}

func (h *chatHandler) GetMessages(c *gin.Context) {
	messages, err := h.domain.GetMessages(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *chatHandler) CreateMessage(c *gin.Context) {
	var body types.CreateMessageParams

	maxFormSize := int64(1<<30) + (1 << 20) // 1gb + 20mb
	if err := c.Request.ParseMultipartForm(maxFormSize); err != nil {
		types.NewAPIError(http.StatusBadRequest, "ERR_PARSE_FORM", "Failed to parse form", err).Respond(c)
		return
	}

	files := c.Request.MultipartForm.File["attachments"]

	if len(files) > 0 {
		if err := validation.ValidateFiles(files, validation.DefaultFileConfig); err != nil {
			err.Respond(c)
			return
		}
	}

	body.ChannelID = c.Request.FormValue("channel_id")
	body.ServerID = c.Request.FormValue("server_id")
	body.MentionsUsers = c.Request.Form["mentions_users"]
	body.MentionsChannels = c.Request.Form["mentions_channels"]
	body.MentionsRoles = c.Request.Form["mentions_roles"]
	contentJSON := c.Request.FormValue("content")
	if err := json.Unmarshal([]byte(contentJSON), &body.Content); err != nil {
		types.NewAPIError(http.StatusBadRequest, "ERR_UNMARSHAL_MESSAGE_CONTENT", "Failed to unmarshal message content.", err).Respond(c)
		return
	}

	if verr := validation.Validate(&body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.CreateMessage(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *chatHandler) EditMessage(c *gin.Context) {
	var body types.EditMessageParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.EditMessage(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *chatHandler) DeleteMessage(c *gin.Context) {
	var body types.DeleteMessageParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.DeleteMessage(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
