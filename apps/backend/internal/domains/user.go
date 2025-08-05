package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/broker"
	"backend/internal/crypto"
	"backend/internal/database"
	"backend/internal/files"
	"backend/internal/types"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetUserByID(ctx *gin.Context, userID string) (*db.User, *types.APIError)
	UpdateEmail(ctx *gin.Context, body *types.UpdateEmailParams) *types.APIError
	UpdateAvatar(ctx *gin.Context, avatar []*multipart.FileHeader, banner []*multipart.FileHeader, body *types.UpdateAvatarParams) (*string, *string, *types.APIError)
	UpdateProfile(ctx *gin.Context, body *types.UpdateProfileParams) *types.APIError
	Setup(ctx *gin.Context) (*types.Setup, *types.APIError)
	UpdatePassword(ctx *gin.Context, body *types.UpdatePasswordParams) *types.APIError
}

type userService struct {
	db     database.Service
	broker broker.Service
	files  files.Service
}

func NewUserService(db database.Service, broker broker.Service, files files.Service) *userService {
	return &userService{
		db:     db,
		broker: broker,
		files:  files,
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

func (s *userService) UpdateEmail(ctx *gin.Context, body *types.UpdateEmailParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}
	userID := u.(*db.User).ID

	updatedUser, err := s.db.UpdateUserEmail(ctx, userID, body)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPDATE_ACCOUNT",
			Cause:   err.Error(),
			Message: "Failed to update account.",
		}
	}

	token, err := ctx.Cookie("token")
	if err != nil {
		return &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_MISSING_TOKEN",
			Message: "Session token not found.",
		}
	}
	s.broker.RefreshCachedUser(ctx, token, updatedUser)

	return nil
}

func (s *userService) UpdateAvatar(ctx *gin.Context, avatar []*multipart.FileHeader, banner []*multipart.FileHeader, body *types.UpdateAvatarParams) (*string, *string, *types.APIError) {
	u, exists := ctx.Get("user")
	if !exists {
		return nil, nil, &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}
	user := u.(*db.User)
	oldAvatar := user.Avatar.String[len(os.Getenv("CDN_URL"))+1:]
	oldBanner := user.Banner.String[len(os.Getenv("CDN_URL"))+1:]

	var avatarURL, bannerURL *string

	if len(avatar) > 0 {
		a, perr := s.files.ProcessAndUploadAvatar(user.ID, "avatar", avatar[0], body.CropAvatar, 1<<20)
		if perr != nil {
			return nil, nil, perr
		}
		avatarURL = a

		err := s.files.DeleteFile(oldAvatar)
		if err != nil {
			fmt.Println("Failed to delete old avatar:", err)
		}
	}

	if len(banner) > 0 {
		b, perr := s.files.ProcessAndUploadAvatar(user.ID, "banner", banner[0], body.CropBanner, 1<<20)
		if perr != nil {
			return nil, nil, perr
		}
		bannerURL = b

		err := s.files.DeleteFile(oldBanner)
		if err != nil {
			fmt.Println("Failed to delete old banner:", err)
		}
	}

	updatedUser, err := s.db.UpdateUserAvatarNBanner(ctx, user.ID, avatarURL, bannerURL)
	if err != nil {
		return nil, nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPDATE_AVATAR_BANNER",
			Cause:   err.Error(),
			Message: "Failed to update avatar/banner.",
		}
	}

	token, err := ctx.Cookie("token")
	if err != nil {
		return nil, nil, &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_MISSING_TOKEN",
			Message: "Session token not found.",
		}
	}
	s.broker.RefreshCachedUser(ctx, token, updatedUser)

	return avatarURL, bannerURL, nil
}

func (s *userService) UpdateProfile(ctx *gin.Context, body *types.UpdateProfileParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}
	userID := u.(*db.User).ID

	var links []types.Link
	var facts []types.Fact

	if err := json.Unmarshal(body.Links, &links); err != nil || len(links) > 2 {
		return &types.APIError{
			Status:  http.StatusBadRequest,
			Code:    "ERR_INVALID_LINKS",
			Message: "You can only have 2 links.",
		}
	}

	if err := json.Unmarshal(body.Facts, &facts); err != nil || len(facts) > 3 {
		return &types.APIError{
			Status:  http.StatusBadRequest,
			Code:    "ERR_INVALID_FACTS",
			Message: "You can only have 3 facts.",
		}
	}

	updatedUser, err := s.db.UpdateUserProfile(ctx, userID, body)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPDATE_PROFILE",
			Cause:   err.Error(),
			Message: "Failed to update profile.",
		}
	}

	token, err := ctx.Cookie("token")
	if err != nil {
		return &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_MISSING_TOKEN",
			Message: "Session token not found.",
		}
	}
	s.broker.RefreshCachedUser(ctx, token, updatedUser)

	return nil
}

func (s *userService) UpdatePassword(ctx *gin.Context, body *types.UpdatePasswordParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}
	userID := u.(*db.User).ID

	password, err := s.db.GetUserPassword(ctx, userID)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_PASSWORD",
			Cause:   err.Error(),
			Message: "Failed to get user password.",
		}
	}

	if valid, err := crypto.VerifyPassword(body.Current, password); err != nil || !valid {
		return &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_INVALID_PASSWORD",
			Message: "Invalid password.",
		}
	}

	hashedPassword, err := crypto.HashPassword(body.New)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_HASH_PASSWORD",
			Cause:   err.Error(),
			Message: "Failed to hash password.",
		}
	}

	if err := s.db.UpdateUserPassword(ctx, userID, hashedPassword); err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPDATE_PASSWORD",
			Cause:   err.Error(),
			Message: "Failed to update password.",
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
