package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/actors"
	"backend/internal/database"
	"backend/internal/permissions"
	"backend/internal/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatService interface {
	CreateMessage(ctx *gin.Context, message *types.CreateMessageParams) *types.APIError
}

type chatService struct {
	db          database.Service
	actors      actors.Service
	permissions permissions.Service
}

func NewChatService(actors actors.Service, db database.Service, permissions permissions.Service) *chatService {
	return &chatService{
		db:          db,
		actors:      actors,
		permissions: permissions,
	}
}

func (s *chatService) CreateMessage(ctx *gin.Context, message *types.CreateMessageParams) *types.APIError {
	user, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized", nil)
	}
	userID := user.(*db.User).ID

	m, err := s.db.CreateMessage(ctx, userID, message)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_CREATE_MESSAGE", "Failed to create message", err)
	}

	fmt.Println(m)

	return nil
}
