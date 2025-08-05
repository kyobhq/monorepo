package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/database"
	"backend/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetUserByID(ctx *gin.Context, userID string) (*db.User, *types.APIError)
	UpdateAccount(ctx *gin.Context, userID string, body *types.UpdateAccountParams) *types.APIError
	UpdateAvatar(ctx *gin.Context, userID string, body *types.UpdateAvatarParams) *types.APIError
	UpdateProfile(ctx *gin.Context, userID string, body *types.UpdateProfileParams) *types.APIError
	Setup(ctx *gin.Context) (*types.Setup, *types.APIError)
}

type userService struct {
	db database.Service
}

func NewUserService(db database.Service) *userService {
	return &userService{
		db: db,
	}
}

func (s *userService) GetUserByID(ctx *gin.Context, userID string) (*db.User, *types.APIError) {
	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_USER_BY_ID",
			Cause:   err.Error(),
			Message: "Failed to get user by ID.",
		}
	}

	return &user, nil
}

func (s *userService) UpdateAccount(ctx *gin.Context, userID string, body *types.UpdateAccountParams) *types.APIError {
	err := s.db.UpdateUserAccount(ctx, userID, body)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPDATE_ACCOUNT",
			Cause:   err.Error(),
			Message: "Failed to update account.",
		}
	}

	return nil
}

func (s *userService) UpdateAvatar(ctx *gin.Context, userID string, body *types.UpdateAvatarParams) *types.APIError {
	err := s.db.UpdateUserAvatarNBanner(ctx, userID, body)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPDATE_AVATAR",
			Cause:   err.Error(),
			Message: "Failed to update avatar.",
		}
	}

	return nil
}

func (s *userService) UpdateProfile(ctx *gin.Context, userID string, body *types.UpdateProfileParams) *types.APIError {
	err := s.db.UpdateUserProfile(ctx, userID, body)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPDATE_PROFILE",
			Cause:   err.Error(),
			Message: "Failed to update profile.",
		}
	}

	return nil
}

func (s *userService) Setup(ctx *gin.Context) (*types.Setup, *types.APIError) {
	var res types.Setup

	u, exists := ctx.Get("user")
	if !exists {
		return nil, &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}
	user := u.(*db.User)

	servers, err := s.db.GetUserServers(ctx, user.ID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_SERVERS",
			Cause:   err.Error(),
			Message: "Failed to get user's servers.",
		}
	}

	res.User = user
	res.Servers = make(map[string]types.ServerWithCategories)
	if len(servers) > 0 {
		serversMap, err := s.processServers(ctx, servers)
		if err != nil {
			return nil, &types.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "ERR_SETUP_SERVERS",
				Cause:   err.Error(),
				Message: "Failed to organize servers for user.",
			}
		}

		res.Servers = serversMap
	}

	return &res, nil
}

func (s *userService) processServers(ctx *gin.Context, servers []db.GetServersFromUserRow) (map[string]types.ServerWithCategories, error) {
	serverIDs := make([]string, 0, len(servers))
	for _, server := range servers {
		serverIDs = append(serverIDs, server.ID)
	}

	allCategories, err := s.db.GetCategoriesFromServers(ctx, serverIDs)
	if err != nil {
		return nil, err
	}

	allChannels, err := s.db.GetChannelsFromServers(ctx, serverIDs)
	if err != nil {
		return nil, err
	}

	allRoles, err := s.db.GetRolesFromServers(ctx, serverIDs)
	if err != nil {
		return nil, err
	}

	categoriesByServer := make(map[string][]db.ChannelCategory)
	for _, category := range allCategories {
		categoriesByServer[category.ServerID] = append(categoriesByServer[category.ServerID], category)
	}

	channelsByCategory := make(map[string][]db.Channel)
	for _, channel := range allChannels {
		channelsByCategory[channel.CategoryID] = append(channelsByCategory[channel.CategoryID], channel)
	}

	rolesByServer := make(map[string][]db.GetRolesFromServersRow)
	for _, role := range allRoles {
		rolesByServer[role.ServerID] = append(rolesByServer[role.ServerID], role)
	}

	result := make(map[string]types.ServerWithCategories)
	for _, server := range servers {
		categoryMap := make(map[string]types.CategoryWithChannels)
		for _, category := range categoriesByServer[server.ID] {
			channelMap := make(map[string]db.Channel)
			for _, channel := range channelsByCategory[category.ID] {
				channelMap[channel.ID] = channel
			}

			categoryMap[category.ID] = types.CategoryWithChannels{
				category,
				channelMap,
			}
		}

		result[server.ID] = types.ServerWithCategories{
			server,
			categoryMap,
			rolesByServer[server.ID],
		}
	}

	return result, nil
}
