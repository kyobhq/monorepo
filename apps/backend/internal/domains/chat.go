package domains

import (
	"backend/internal/actors"
	"backend/internal/database"
	"backend/internal/permissions"
	"backend/internal/types"

	"github.com/gin-gonic/gin"
)

type ChatService interface {
	CreateMessage(ctx *gin.Context, message *types.ChatMessage) *types.APIError
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

func (s *chatService) CreateMessage(ctx *gin.Context, message *types.ChatMessage) *types.APIError {
	//TODO: send to channel actor
	// protoMessage := &proto.NewChatMessage{
	// 	Message: &proto.ChatMessage{
	// 		AuthorId: string,
	//
	// 	},
	// }
	//
	// userPID := ch.actors.GetUser(body.AuthorID)
	// channelPID := ch.actors.GetChannel(body.ServerID, body.ChannelID)
	// ch.actors.SendMessageTo(actors.ServerEngine, channelPID, protoMessage, userPID)
	return nil
}
