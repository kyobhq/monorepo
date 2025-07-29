package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/database"
	"backend/internal/files"
	"backend/internal/types"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServerService interface {
	CreateServer(ctx *gin.Context, serverAvatar []*multipart.FileHeader, body *types.CreateServerParams) (*db.Server, *types.APIError)
	JoinServer(ctx *gin.Context, serverID string) (*db.Server, *types.APIError)
	LeaveServer(ctx *gin.Context, serverID string) *types.APIError
	CreateInvite(ctx *gin.Context, serverID string) (*string, *types.APIError)
	DeleteInvite(ctx *gin.Context, inviteID string) *types.APIError
	EditProfile(ctx *gin.Context, serverID string, body *types.UpdateServerProfileParams) *types.APIError
	EditAvatar(ctx *gin.Context, serverID string, body *types.UpdateServerAvatarParams) *types.APIError
	DeleteServer(ctx *gin.Context, serverID string) *types.APIError
}

type serverService struct {
	db    database.Service
	files files.Service
}

func NewServerService(db database.Service, files files.Service) *serverService {
	return &serverService{
		db:    db,
		files: files,
	}
}

func (s *serverService) CreateServer(ctx *gin.Context, serverAvatar []*multipart.FileHeader, body *types.CreateServerParams) (*db.Server, *types.APIError) {
	user, exists := ctx.Get("user")
	if !exists {
		return nil, &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}

	avatarURL, perr := s.files.ProcessAndUploadImage(serverAvatar[0], body.Crop, 1<<20)
	if perr != nil {
		return nil, perr
	}

	userID := user.(*db.User).ID

	server, err := s.db.CreateServer(ctx, userID, body, avatarURL)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CREATE_SERVER",
			Cause:   err.Error(),
			Message: "Failed to create server.",
		}
	}

	return server, nil
}

func (s *serverService) JoinServer(ctx *gin.Context, serverID string) (*db.Server, *types.APIError) {
	user, exists := ctx.Get("user")
	if !exists {
		return nil, &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}

	server, err := s.db.GetServer(ctx, serverID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_SERVER",
			Cause:   err.Error(),
			Message: "Failed to get server.",
		}
	}

	userID := user.(*db.User).ID

	err = s.db.JoinServer(ctx, server.ID, userID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_JOIN_SERVER",
			Cause:   err.Error(),
			Message: "Failed to join server.",
		}
	}

	return &server, nil
}

func (s *serverService) LeaveServer(ctx *gin.Context, serverID string) *types.APIError {
	user, exists := ctx.Get("user")
	if !exists {
		return &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}

	userID := user.(db.User).ID

	err := s.db.LeaveServer(ctx, serverID, userID)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_LEAVE_SERVER",
			Cause:   err.Error(),
			Message: "Failed to leave server.",
		}
	}

	return nil
}

func (s *serverService) CreateInvite(ctx *gin.Context, serverID string) (*string, *types.APIError) {
	inviteID, err := s.db.CreateInvite(ctx, serverID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CREATE_INVITE",
			Cause:   err.Error(),
			Message: "Failed to create invite.",
		}
	}

	inviteURL := fmt.Sprintf("%s/invite/%s", "http://localhost:5173", inviteID)

	return &inviteURL, nil
}

func (s *serverService) EditProfile(ctx *gin.Context, serverID string, body *types.UpdateServerProfileParams) *types.APIError {
	err := s.db.UpdateServerProfile(ctx, serverID, body)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPDATE_SERVER_PROFILE",
			Cause:   err.Error(),
			Message: "Failed to update server profile.",
		}
	}

	return nil
}

func (s *serverService) DeleteServer(ctx *gin.Context, serverID string) *types.APIError {
	pgconn, err := s.db.DeleteServer(ctx, serverID)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_DELETE_SERVER",
			Cause:   err.Error(),
			Message: "Failed to delete server.",
		}
	}

	if pgconn.RowsAffected() == 0 {
		return &types.APIError{
			Status:  http.StatusNotFound,
			Code:    "ERR_SERVER_NOT_FOUND",
			Message: "Server not found.",
		}
	}

	return nil
}

func (s *serverService) DeleteInvite(ctx *gin.Context, inviteID string) *types.APIError {
	return nil
}

func (s *serverService) EditAvatar(ctx *gin.Context, serverID string, body *types.UpdateServerAvatarParams) *types.APIError {
	return nil
}
