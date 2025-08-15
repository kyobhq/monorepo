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
	"github.com/nrednav/cuid2"
)

type UserService interface {
	GetUserByID(ctx *gin.Context, userID string) (*db.User, *types.APIError)
	UpdateEmail(ctx *gin.Context, body *types.UpdateEmailParams) *types.APIError
	UpdateAvatar(ctx *gin.Context, avatar []*multipart.FileHeader, banner []*multipart.FileHeader, body *types.UpdateAvatarParams) (*string, *string, *types.APIError)
	UpdateProfile(ctx *gin.Context, body *types.UpdateProfileParams) *types.APIError
	Setup(ctx *gin.Context) (*types.Setup, *types.APIError)
	UpdatePassword(ctx *gin.Context, body *types.UpdatePasswordParams) *types.APIError
	GetUserProfile(ctx *gin.Context, userID string) (*db.GetUserProfileRow, *types.APIError)
	UploadEmojis(ctx *gin.Context, emojis []*multipart.FileHeader, shortcodes []string, body *types.UploadEmojiParams) (*[]types.EmojiResponse, *types.APIError)
	UpdateEmoji(ctx *gin.Context, body *types.UpdateEmojiParams) *types.APIError
	DeleteEmoji(ctx *gin.Context) *types.APIError
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
	var oldAvatar, oldBanner string

	if user.Avatar.String != "" {
		oldAvatar = user.Avatar.String[len(os.Getenv("CDN_URL"))+1:]
	}
	if user.Banner.String != "" {
		oldBanner = user.Banner.String[len(os.Getenv("CDN_URL"))+1:]
	}

	var avatarURL, bannerURL *string

	if len(avatar) > 0 {
		a, perr := s.files.ProcessAndUploadAvatar(user.ID, "avatar", avatar[0], body.CropAvatar)
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
		b, perr := s.files.ProcessAndUploadAvatar(user.ID, "banner", banner[0], body.CropBanner)
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
	user.Password = ""

	emojis, err := s.db.GetEmojis(ctx, user.ID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_EMOJIS",
			Cause:   err.Error(),
			Message: "Failed to get user's emojis.",
		}
	}

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
	res.Emojis = emojis
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
	for _, server := range servers {
		rolesByServer[server.ID] = []db.GetRolesFromServersRow{}
	}
	for _, role := range allRoles {
		rolesByServer[role.ServerID] = append(rolesByServer[role.ServerID], role)
	}

	result := make(map[string]types.ServerWithCategories)
	members := make([]types.Member, 0)
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
			members,
			rolesByServer[server.ID],
		}
	}

	return result, nil
}

func (s *userService) GetUserProfile(ctx *gin.Context, userID string) (*db.GetUserProfileRow, *types.APIError) {
	user, err := s.db.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_GET_USER_PROFILE",
			Cause:   err.Error(),
			Message: "Failed to get user profile.",
		}
	}

	return &user, nil
}

func (s *userService) UploadEmojis(ctx *gin.Context, emojis []*multipart.FileHeader, shortcodes []string, body *types.UploadEmojiParams) (*[]types.EmojiResponse, *types.APIError) {
	u, exists := ctx.Get("user")
	if !exists {
		return nil, &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}
	userID := u.(*db.User).ID

	emojisURLs, err := s.files.ProcessAndUploadEmojis(emojis)
	if err != nil {
		return nil, err
	}

	var emojisToUpload []db.CreateEmojiParams
	var emojisResponse []types.EmojiResponse

	for i, url := range emojisURLs {
		emojiID := cuid2.Generate()
		emojisToUpload = append(emojisToUpload, db.CreateEmojiParams{
			ID:        emojiID,
			UserID:    userID,
			Shortcode: body.Shortcodes[i],
			Url:       url,
		})

		emojisResponse = append(emojisResponse, types.EmojiResponse{
			ID:        emojiID,
			Shortcode: body.Shortcodes[i],
			URL:       url,
		})
	}

	if err := s.db.UploadEmojis(ctx, userID, emojisToUpload); err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPLOAD_EMOJIS",
			Cause:   err.Error(),
			Message: "Failed to upload emojis.",
		}
	}

	return &emojisResponse, nil
}

func (s *userService) UpdateEmoji(ctx *gin.Context, body *types.UpdateEmojiParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}

	userID := u.(*db.User).ID
	emojiID := ctx.Param("emoji_id")

	if err := s.db.UpdateEmoji(ctx, emojiID, userID, body); err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPDATE_EMOJI",
			Cause:   err.Error(),
			Message: "Failed to update emoji.",
		}
	}

	return nil
}

func (s *userService) DeleteEmoji(ctx *gin.Context) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_UNAUTHORIZED",
			Message: "Unauthorized.",
		}
	}

	userID := u.(*db.User).ID
	emojiID := ctx.Param("emoji_id")

	if err := s.db.DeleteEmoji(ctx, emojiID, userID); err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_DELETE_EMOJI",
			Cause:   err.Error(),
			Message: "Failed to delete emoji.",
		}
	}

	return nil
}
