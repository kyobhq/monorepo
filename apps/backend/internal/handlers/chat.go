package handlers

import (
	"backend/internal/domains"
	"backend/internal/types"
	"backend/internal/validation"
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

func (ch *chatHandler) CreateMessage(c *gin.Context) {
	var body types.ChatMessage

	maxFormSize := int64(1<<30) + (1 << 20)                           // 1gb + 20mb
	if err := c.Request.ParseMultipartForm(maxFormSize); err != nil { // 32MB max
		types.NewAPIError(http.StatusBadRequest, "ERR_PARSE_FORM", "Failed to parse form", err).Respond(c)
		return
	}

	files := c.Request.MultipartForm.File["attachments"]

	if err := validation.ValidateFiles(files, validation.DefaultFileConfig); err != nil {
		err.Respond(c)
		return
	}

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := ch.domain.CreateMessage(&body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
