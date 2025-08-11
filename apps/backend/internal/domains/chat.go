package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/actors"
	"backend/internal/database"
	"backend/internal/files"
	"backend/internal/permissions"
	"backend/internal/types"
	"backend/proto"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatService interface {
	CreateMessage(ctx *gin.Context, files []*multipart.FileHeader, message *types.CreateMessageParams) *types.APIError
	EditMessage(ctx *gin.Context, message *types.EditMessageParams) *types.APIError
	DeleteMessage(ctx *gin.Context, params *types.DeleteMessageParams) *types.APIError
	GetMessages(ctx *gin.Context) ([]db.GetMessagesFromChannelRow, *types.APIError)
}

type chatService struct {
	db          database.Service
	actors      actors.Service
	permissions permissions.Service
	files       files.Service
}

func NewChatService(actors actors.Service, db database.Service, files files.Service, permissions permissions.Service) *chatService {
	return &chatService{
		db:          db,
		actors:      actors,
		permissions: permissions,
		files:       files,
	}
}

func (s *chatService) GetMessages(ctx *gin.Context) ([]db.GetMessagesFromChannelRow, *types.APIError) {
	channelID := ctx.Param("channel_id")

	messages, err := s.db.GetMessages(ctx, channelID)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_GET_MESSAGES", "Failed to get messages", err)
	}

	return messages, nil
}

func (s *chatService) CreateMessage(ctx *gin.Context, files []*multipart.FileHeader, message *types.CreateMessageParams) *types.APIError {
	user, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized", nil)
	}
	author := user.(*db.User)

	jsonAttachments, ferr := s.files.ProcessAndUploadFiles(files)
	if ferr != nil {
		return ferr
	}

	message.Attachments = jsonAttachments

	m, err := s.db.CreateMessage(ctx, author.ID, message)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_CREATE_MESSAGE", "Failed to create message", err)
	}

	pbMessage := &proto.NewChatMessage{
		Message: &proto.Message{
			Id: m.ID,
			Author: &proto.User{
				Id:          author.ID,
				Avatar:      author.Avatar.String,
				DisplayName: author.DisplayName,
			},
			ServerId:         m.ServerID,
			ChannelId:        m.ChannelID,
			Content:          m.Content,
			Everyone:         m.Everyone,
			MentionsUsers:    m.MentionsUsers,
			MentionsChannels: m.MentionsChannels,
			Attachments:      m.Attachments,
			CreatedAt:        timestamppb.New(m.CreatedAt),
			UpdatedAt:        timestamppb.New(m.UpdatedAt),
		},
	}

	s.actors.SendChatMessage(pbMessage)

	return nil
}

func (s *chatService) EditMessage(ctx *gin.Context, message *types.EditMessageParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized", nil)
	}
	userID := u.(*db.User).ID
	messageID := ctx.Param("message_id")
	authorID, err := s.db.GetMessageAuthor(ctx, messageID)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_GET_AUTHOR", "Failed to get message author", nil)
	}

	if authorID != userID {
		return types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN", "You are not allowed to edit this message", nil)
	}

	err = s.db.EditMessage(ctx, messageID, message)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_EDIT_MESSAGE", "Failed to edit message", err)
	}

	s.actors.EditMessage(&proto.EditChatMessage{
		Message: &proto.Message{
			Id:               messageID,
			ServerId:         message.ServerID,
			ChannelId:        message.ChannelID,
			Content:          message.Content,
			Everyone:         message.Everyone,
			MentionsUsers:    message.MentionsUsers,
			MentionsChannels: message.MentionsChannels,
			UpdatedAt:        timestamppb.New(time.Now()),
		},
	})

	return nil
}

func (s *chatService) DeleteMessage(ctx *gin.Context, params *types.DeleteMessageParams) *types.APIError {
	messageID := ctx.Param("message_id")

	if allowed := s.permissions.CheckPermission(ctx, params.ServerID, types.ManageMessages, messageID, params.AuthorID); !allowed {
		return &types.APIError{
			Status:  http.StatusForbidden,
			Code:    "ERR_FORBIDDEN",
			Message: "You are not allowed to delete this message.",
			Cause:   "",
		}
	}

	err := s.db.DeleteMessage(ctx, messageID, params.AuthorID)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_DELETE_MESSAGE",
			Message: "Failed to delete message.",
			Cause:   err.Error(),
		}
	}

	s.actors.DeleteMessage(&proto.DeleteChatMessage{
		Message: &proto.Message{
			Id:        messageID,
			ServerId:  params.ServerID,
			ChannelId: params.ChannelID,
		},
	})

	return nil
}
