package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/actors"
	"backend/internal/database"
	"backend/internal/files"
	"backend/internal/permissions"
	"backend/internal/types"
	"backend/proto"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ServerService interface {
	CreateServer(ctx *gin.Context, serverAvatar []*multipart.FileHeader, body *types.CreateServerParams) (*db.Server, *types.APIError)
	JoinServer(ctx *gin.Context, body *types.JoinServerParams) (*types.JoinServerWithCategories, *types.APIError)
	LeaveServer(ctx *gin.Context, serverID string) *types.APIError
	CreateInvite(ctx *gin.Context, serverID string) (*string, *types.APIError)
	DeleteInvite(ctx *gin.Context, inviteID string) *types.APIError
	EditProfile(ctx *gin.Context, serverID string, body *types.UpdateServerProfileParams) *types.APIError
	EditAvatar(ctx *gin.Context, serverID string, body *types.UpdateServerAvatarParams) *types.APIError
	DeleteServer(ctx *gin.Context, serverID string) *types.APIError
	GetInformations(ctx *gin.Context) (*db.GetServerInformationsRow, *types.APIError)
}

type serverService struct {
	db          database.Service
	files       files.Service
	actors      actors.Service
	permissions permissions.Service
}

func NewServerService(db database.Service, actors actors.Service, files files.Service, permissions permissions.Service) *serverService {
	return &serverService{
		db:          db,
		files:       files,
		actors:      actors,
		permissions: permissions,
	}
}

func (s *serverService) CreateServer(ctx *gin.Context, serverAvatar []*multipart.FileHeader, body *types.CreateServerParams) (*db.Server, *types.APIError) {
	u, exists := ctx.Get("user")
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

	user := u.(*db.User)

	server, err := s.db.CreateServer(ctx, user.ID, body, avatarURL)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CREATE_SERVER",
			Cause:   err.Error(),
			Message: "Failed to create server.",
		}
	}

	s.actors.StartServerInRegion(server.ID, os.Getenv("REGION"))
	userPID := s.actors.GetUser(user.ID)

	changeStatus := &proto.ChangeStatus{
		Type: "connect",
		User: &proto.User{
			Id:          user.ID,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar.String,
		},
		ServerId: server.ID,
		Status:   "online",
	}
	s.actors.SendUserStatusMessage(userPID, changeStatus)

	return server, nil
}

func (s *serverService) JoinServer(ctx *gin.Context, body *types.JoinServerParams) (*types.JoinServerWithCategories, *types.APIError) {
	u, exists := ctx.Get("user")
	if !exists {
		return nil, &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}
	user := u.(*db.User)

	var serverID string

	if body.InviteID != "" {
		id, err := s.db.CheckInvite(ctx, body.InviteID)
		if err != nil {
			return nil, &types.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "ERR_CHECK_INVITE",
				Cause:   err.Error(),
				Message: "Failed to check invite.",
			}
		}
		serverID = id
	}

	if body.ServerID != "" {
		server, err := s.db.GetServer(ctx, body.ServerID)
		if err != nil {
			return nil, &types.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "ERR_GET_SERVER",
				Cause:   err.Error(),
				Message: "Failed to get server.",
			}
		}
		serverID = server.ID
	}

	server, err := s.db.JoinServer(ctx, serverID, user.ID, body.Position)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_JOIN_SERVER",
			Cause:   err.Error(),
			Message: "Failed to join server.",
		}
	}

	userPID := s.actors.GetUser(user.ID)
	changeStatus := &proto.ChangeStatus{
		Type: "join",
		User: &proto.User{
			Id:          user.ID,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar.String,
		},
		ServerId: server.ID,
		Status:   "online",
	}
	s.actors.SendUserStatusMessage(userPID, changeStatus)

	categories, err := s.db.GetCategoriesFromServer(ctx, serverID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_CATEGORIES",
			Cause:   err.Error(),
			Message: "Failed to get categories.",
		}
	}

	channels, err := s.db.GetChannelsFromServer(ctx, serverID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_CHANNELS",
			Cause:   err.Error(),
			Message: "Failed to get channels.",
		}
	}

	roles, err := s.db.GetRolesFromServer(ctx, serverID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_ROLES",
			Cause:   err.Error(),
			Message: "Failed to get roles.",
		}
	}

	categoryMap := make(map[string]types.CategoryWithChannels)
	for _, category := range categories {
		channelMap := make(map[string]db.Channel)
		for _, channel := range channels {
			if channel.CategoryID == category.ID {
				channelMap[channel.ID] = channel
			}
		}

		categoryMap[category.ID] = types.CategoryWithChannels{
			category,
			channelMap,
		}
	}

	serverWithCategories := &types.JoinServerWithCategories{
		server,
		categoryMap,
		roles,
	}

	return serverWithCategories, nil
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

	userID := user.(*db.User).ID

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
	user, exists := ctx.Get("user")
	if !exists {
		return nil, &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}
	userID := user.(*db.User).ID

	inviteID, err := s.db.CreateInvite(ctx, userID, serverID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CREATE_INVITE",
			Cause:   err.Error(),
			Message: "Failed to create invite.",
		}
	}

	inviteURL := fmt.Sprintf("%s/invite/%s", os.Getenv("DOMAIN"), inviteID)

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

func (s *serverService) GetInformations(ctx *gin.Context) (*db.GetServerInformationsRow, *types.APIError) {
	serverID := ctx.Param("server_id")

	userIDs := s.actors.GetActiveUsers(serverID)
	serverInformations, err := s.db.GetServerInformations(ctx, serverID, userIDs)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_SERVER_INFORMATIONS",
			Cause:   err.Error(),
			Message: "Failed to get server informations.",
		}
	}

	return &serverInformations, nil
}

func (s *serverService) DeleteInvite(ctx *gin.Context, inviteID string) *types.APIError {
	return nil
}

func (s *serverService) EditAvatar(ctx *gin.Context, serverID string, body *types.UpdateServerAvatarParams) *types.APIError {
	return nil
}
