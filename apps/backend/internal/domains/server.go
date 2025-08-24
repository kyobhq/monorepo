package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/actors"
	"backend/internal/database"
	"backend/internal/files"
	"backend/internal/permissions"
	"backend/internal/types"
	"backend/internal/validation"
	"backend/proto"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nrednav/cuid2"
)

type ServerService interface {
	CreateServer(ctx *gin.Context, serverAvatar []*multipart.FileHeader, body *types.CreateServerParams) (*db.Server, *types.APIError)
	JoinServer(ctx *gin.Context, body *types.JoinServerParams) (*types.JoinServerWithCategories, *types.APIError)
	LeaveServer(ctx *gin.Context) *types.APIError
	CreateInvite(ctx *gin.Context) (*string, *types.APIError)
	DeleteInvite(ctx *gin.Context) *types.APIError
	UpdateProfile(ctx *gin.Context, body *types.UpdateServerProfileParams) *types.APIError
	UpdateAvatar(ctx *gin.Context, avatar []*multipart.FileHeader, banner []*multipart.FileHeader, body *types.UpdateAvatarParams) (*string, *string, *types.APIError)
	DeleteServer(ctx *gin.Context) *types.APIError
	GetInformations(ctx *gin.Context) (*db.GetServerInformationsRow, *types.APIError)
	GetMembers(ctx *gin.Context) ([]db.GetServerMembersRow, *types.APIError)
	GetBannedMembers(ctx *gin.Context) ([]db.GetBannedMembersRow, *types.APIError)
	BanUser(ctx *gin.Context, body *types.BanUserParams) *types.APIError
	UnbanUser(ctx *gin.Context) *types.APIError
	KickUser(ctx *gin.Context, body *types.KickUserParams) *types.APIError
	SearchMembers(ctx *gin.Context) ([]db.SearchServerMembersRow, *types.APIError)
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
		return nil, types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}

	avatarURL, perr := s.files.ProcessAndUploadAvatar(cuid2.Generate(), "avatar", serverAvatar[0], body.Crop)
	if perr != nil {
		return nil, perr
	}

	user := u.(*db.User)

	server, err := s.db.CreateServer(ctx, user.ID, body, avatarURL)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_CREATE_SERVER", "Failed to create server.", err)
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
		return nil, types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	user := u.(*db.User)

	var serverID string
	if body.InviteID != "" {
		id, err := s.db.CheckInvite(ctx, body.InviteID)
		if err != nil {
			return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_CHECK_INVITE", "Failed to check invite.", err)
		}
		serverID = id
	}

	reason, err := s.db.CheckBan(ctx, serverID, user.ID)
	if err == nil {
		return nil, types.NewAPIError(http.StatusForbidden, "USER_BANNED", reason.String, nil)
	}

	server, categories, channels, roles, latestMessagesSent, err := s.db.JoinServer(ctx, serverID, user.ID, body.Position)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_JOIN_SERVER", "Failed to join server.", err)
	}

	allMessagesSentMap := make(map[string]string)
	for _, message := range latestMessagesSent {
		allMessagesSentMap[message.ChannelID] = message.ID
	}

	categoryMap := make(map[string]types.CategoryWithChannels)
	for _, category := range categories {
		channelMap := make(map[string]types.ServerChannel)
		for _, channel := range channels {
			if channel.CategoryID.String == category.ID {
				channelMap[channel.ID] = types.ServerChannel{
					channel,
					allMessagesSentMap[channel.ID],
					allMessagesSentMap[channel.ID],
					json.RawMessage(`[]`),
				}
			}
		}

		categoryMap[category.ID] = types.CategoryWithChannels{
			category,
			channelMap,
		}
	}

	serverWithCategories := &types.JoinServerWithCategories{
		*server,
		categoryMap,
		roles,
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

	return serverWithCategories, nil
}

func (s *serverService) LeaveServer(ctx *gin.Context) *types.APIError {
	user, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}

	userID := user.(*db.User).ID
	serverID := ctx.Param("server_id")

	err := s.db.LeaveServer(ctx, serverID, userID)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_LEAVE_SERVER", "Failed to leave server.", err)
	}

	s.actors.LeaveServer(serverID, userID)

	return nil
}

func (s *serverService) CreateInvite(ctx *gin.Context) (*string, *types.APIError) {
	user, exists := ctx.Get("user")
	if !exists {
		return nil, types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}

	userID := user.(*db.User).ID
	serverID := ctx.Param("server_id")

	inviteID, err := s.db.CreateInvite(ctx, userID, serverID)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_CREATE_INVITE", "Failed to create invite.", err)
	}

	inviteURL := fmt.Sprintf("https://%s/invite/%s", os.Getenv("DOMAIN"), inviteID)

	return &inviteURL, nil
}

func (s *serverService) UpdateProfile(ctx *gin.Context, body *types.UpdateServerProfileParams) *types.APIError {
	serverID := ctx.Param("server_id")

	if allowed := s.permissions.CheckPermission(ctx, serverID, types.ManageServer); !allowed {
		return types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN", "You are not allowed to edit this server.", nil)
	}

	err := s.db.UpdateServerProfile(ctx, serverID, body)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_UPDATE_SERVER_PROFILE", "Failed to update server profile.", err)
	}

	s.actors.ProfileServerChange(serverID, body)

	return nil
}

func (s *serverService) DeleteServer(ctx *gin.Context) *types.APIError {
	user, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	userID := user.(*db.User).ID
	serverID := ctx.Param("server_id")

	err := s.db.DeleteServer(ctx, userID, serverID)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_DELETE_SERVER", "Failed to delete server.", err)
	}

	s.actors.KillServer(serverID)

	return nil
}

func (s *serverService) GetInformations(ctx *gin.Context) (*db.GetServerInformationsRow, *types.APIError) {
	user, exists := ctx.Get("user")
	if !exists {
		return nil, types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	userID := user.(*db.User).ID
	serverID := ctx.Param("server_id")

	userIDs := s.actors.GetActiveUsers(serverID)
	serverInformations, err := s.db.GetServerInformations(ctx, userID, serverID, userIDs)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_GET_SERVER_INFORMATIONS", "Failed to get server informations.", err)
	}

	return &serverInformations, nil
}

func (s *serverService) GetMembers(ctx *gin.Context) ([]db.GetServerMembersRow, *types.APIError) {
	serverID := ctx.Param("server_id")

	offsetStr := ctx.DefaultQuery("offset", "0")
	offset := 0
	if o, err := fmt.Sscanf(offsetStr, "%d", &offset); err != nil || o != 1 {
		offset = 0
	}

	userIDs := s.actors.GetActiveUsers(serverID)
	members, err := s.db.GetServerMembers(ctx, serverID, int32(offset), userIDs)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_GET_MEMBERS", "Failed to get server members.", err)
	}

	return members, nil
}

func (s *serverService) SearchMembers(ctx *gin.Context) ([]db.SearchServerMembersRow, *types.APIError) {
	serverID := ctx.Param("server_id")
	query := validation.SanitizeQuery(ctx.Query("query"))

	members, err := s.db.SearchServerMembers(ctx, serverID, query)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_SEARCH_MEMBERS", "Failed to search members.", err)
	}

	return members, nil
}

func (s *serverService) DeleteInvite(ctx *gin.Context) *types.APIError {
	return nil
}

func (s *serverService) UpdateAvatar(ctx *gin.Context, avatar []*multipart.FileHeader, banner []*multipart.FileHeader, body *types.UpdateAvatarParams) (*string, *string, *types.APIError) {
	serverID := ctx.Param("server_id")

	if allowed := s.permissions.CheckPermission(ctx, serverID, types.ManageServer); !allowed {
		return nil, nil, types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN", "You are not allowed to edit this server.", nil)
	}

	server, err := s.db.GetServer(ctx, serverID)
	if err != nil {
		return nil, nil, types.NewAPIError(http.StatusInternalServerError, "ERR_GET_SERVER", "Failed to get server.", err)
	}

	var oldAvatar, oldBanner string

	if server.Avatar.String != "" {
		oldAvatar = server.Avatar.String[len(os.Getenv("CDN_URL"))+1:]
	}
	if server.Banner.String != "" {
		oldBanner = server.Banner.String[len(os.Getenv("CDN_URL"))+1:]
	}

	var avatarURL, bannerURL *string

	if len(avatar) > 0 {
		a, perr := s.files.ProcessAndUploadAvatar(server.ID, "avatar", avatar[0], body.CropAvatar)
		if perr != nil {
			return nil, nil, perr
		}
		avatarURL = a

		if oldAvatar != "" {
			err := s.files.DeleteFile(oldAvatar)
			if err != nil {
				fmt.Println("Failed to delete old avatar:", err)
			}
		}
	}

	if len(banner) > 0 {
		b, perr := s.files.ProcessAndUploadAvatar(server.ID, "banner", banner[0], body.CropBanner)
		if perr != nil {
			return nil, nil, perr
		}
		bannerURL = b

		if oldBanner != "" {
			err := s.files.DeleteFile(oldBanner)
			if err != nil {
				fmt.Println("Failed to delete old banner:", err)
			}
		}
	}

	err = s.db.UpdateServerAvatarNBanner(ctx, serverID, avatarURL, bannerURL)
	if err != nil {
		return nil, nil, types.NewAPIError(http.StatusInternalServerError, "ERR_UPDATE_SERVER_AVATAR", "Failed to update server avatar/banner.", err)
	}

	s.actors.AvatarServerChange(serverID, bannerURL, avatarURL)

	return avatarURL, bannerURL, nil
}

func (s *serverService) BanUser(ctx *gin.Context, body *types.BanUserParams) *types.APIError {
	serverID := ctx.Param("server_id")
	if allowed := s.permissions.CheckPermission(ctx, serverID, types.BanMembers); !allowed {
		return types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN", "You are not allowed to ban users.", nil)
	}

	if err := s.db.BanUser(ctx, serverID, body); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_BAN_USER", "Failed to ban user.", err)
	}

	s.actors.BanUser(serverID, body)

	return nil
}

func (s *serverService) UnbanUser(ctx *gin.Context) *types.APIError {
	serverID := ctx.Param("server_id")
	if allowed := s.permissions.CheckPermission(ctx, serverID, types.BanMembers); !allowed {
		return types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN", "You are not allowed to unban users.", nil)
	}

	userID := ctx.Param("user_id")
	if err := s.db.UnbanUser(ctx, serverID, userID); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_UNBAN_USER", "Failed to unban user.", err)
	}

	return nil
}

func (s *serverService) KickUser(ctx *gin.Context, body *types.KickUserParams) *types.APIError {
	serverID := ctx.Param("server_id")

	if allowed := s.permissions.CheckPermission(ctx, serverID, types.KickMembers); !allowed {
		return types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN", "You are not allowed to kick users.", nil)
	}

	if err := s.db.KickUser(ctx, serverID, body); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_KICK_USER", "Failed to kick user.", err)
	}

	s.actors.KickUser(serverID, body)

	return nil
}

func (s *serverService) GetBannedMembers(ctx *gin.Context) ([]db.GetBannedMembersRow, *types.APIError) {
	serverID := ctx.Param("server_id")

	if allowed := s.permissions.CheckPermission(ctx, serverID, types.BanMembers); !allowed {
		return nil, types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN", "You are not allowed to get banned users.", nil)
	}

	bans, err := s.db.GetBannedMembers(ctx, serverID)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_GET_BANNED_MEMBERS", "Failed to get banned members.", err)
	}

	return bans, nil
}
