package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/actors"
	"backend/internal/broker"
	"backend/internal/crypto"
	"backend/internal/database"
	"backend/internal/files"
	"backend/internal/types"
	"backend/internal/utils"
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
	GetUserProfile(ctx *gin.Context) (*db.GetUserProfileRow, *types.APIError)
	UploadEmojis(ctx *gin.Context, emojis []*multipart.FileHeader, shortcodes []string, body *types.UploadEmojiParams) (*[]types.EmojiResponse, *types.APIError)
	UpdateEmoji(ctx *gin.Context, body *types.UpdateEmojiParams) *types.APIError
	DeleteEmoji(ctx *gin.Context) *types.APIError
	DeleteAccount(ctx *gin.Context) *types.APIError
	Sync(ctx *gin.Context, body *types.SyncParams) *types.APIError
}

type userService struct {
	db     database.Service
	broker broker.Service
	files  files.Service
	actors actors.Service
}

func NewUserService(db database.Service, broker broker.Service, files files.Service, actors actors.Service) *userService {
	return &userService{
		db:     db,
		broker: broker,
		files:  files,
		actors: actors,
	}
}

func (s *userService) GetUserByID(ctx *gin.Context, userID string) (*db.User, *types.APIError) {
	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_GET_USER_BY_ID", "Failed to get user by ID.", err)
	}

	return &user, nil
}

func (s *userService) UpdateEmail(ctx *gin.Context, body *types.UpdateEmailParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	userID := u.(*db.User).ID

	updatedUser, err := s.db.UpdateUserEmail(ctx, userID, body)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_UPDATE_ACCOUNT", "Failed to update account.", err)
	}

	token, err := ctx.Cookie("token")
	if err != nil {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_MISSING_TOKEN", "Session token not found.", err)
	}
	s.broker.RefreshCachedUser(ctx, token, updatedUser)

	return nil
}

func (s *userService) UpdateAvatar(ctx *gin.Context, avatar []*multipart.FileHeader, banner []*multipart.FileHeader, body *types.UpdateAvatarParams) (*string, *string, *types.APIError) {
	u, exists := ctx.Get("user")
	if !exists {
		return nil, nil, types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
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
		return nil, nil, types.NewAPIError(http.StatusInternalServerError, "ERR_UPDATE_AVATAR_BANNER", "Failed to update avatar/banner.", err)
	}

	token, err := ctx.Cookie("token")
	if err != nil {
		return nil, nil, types.NewAPIError(http.StatusUnauthorized, "ERR_MISSING_TOKEN", "Session token not found.", err)
	}
	s.broker.RefreshCachedUser(ctx, token, updatedUser)

	return avatarURL, bannerURL, nil
}

func (s *userService) UpdateProfile(ctx *gin.Context, body *types.UpdateProfileParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	userID := u.(*db.User).ID

	var links []types.Link
	var facts []types.Fact

	if err := json.Unmarshal(body.Links, &links); err != nil || len(links) > 2 {
		return types.NewAPIError(http.StatusBadRequest, "ERR_INVALID_LINKS", "You can only have 2 links.", nil)
	}

	if err := json.Unmarshal(body.Facts, &facts); err != nil || len(facts) > 3 {
		return types.NewAPIError(http.StatusBadRequest, "ERR_INVALID_FACTS", "You can only have 3 facts.", nil)
	}

	updatedUser, err := s.db.UpdateUserProfile(ctx, userID, body)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_UPDATE_PROFILE", "Failed to update profile.", err)
	}

	token, err := ctx.Cookie("token")
	if err != nil {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_MISSING_TOKEN", "Session token not found.", err)
	}
	s.broker.RefreshCachedUser(ctx, token, updatedUser)

	return nil
}

func (s *userService) UpdatePassword(ctx *gin.Context, body *types.UpdatePasswordParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	userID := u.(*db.User).ID

	password, err := s.db.GetUserPassword(ctx, userID)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_GET_PASSWORD", "Failed to get user password.", err)
	}

	if valid, err := crypto.VerifyPassword(body.Current, password); err != nil || !valid {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_INVALID_PASSWORD", "Invalid password.", err)
	}

	hashedPassword, err := crypto.HashPassword(body.New)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_HASH_PASSWORD", "Failed to hash password.", err)
	}

	if err := s.db.UpdateUserPassword(ctx, userID, hashedPassword); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_UPDATE_PASSWORD", "Failed to update password.", err)
	}

	return nil
}

type messageStateMaps struct {
	readMap     map[string]string
	sentMap     map[string]string
	mentionsMap map[string]json.RawMessage
}

func (s *userService) buildMessageStateMaps(ctx *gin.Context, userID string, channelIDs []string) (*messageStateMaps, error) {
	if len(channelIDs) == 0 {
		return &messageStateMaps{
			readMap:     make(map[string]string),
			sentMap:     make(map[string]string),
			mentionsMap: make(map[string]json.RawMessage),
		}, nil
	}

	messagesRead, err := s.db.GetLatestMessagesRead(ctx, userID)
	if err != nil {
		return nil, err
	}

	messagesSent, err := s.db.GetLatestMessagesSent(ctx, channelIDs)
	if err != nil {
		return nil, err
	}

	readMap := make(map[string]string)
	mentionsMap := make(map[string]json.RawMessage)
	sentMap := make(map[string]string)

	for _, message := range messagesRead {
		readMap[message.ChannelID] = message.LastReadMessageID.String
		mentionsMap[message.ChannelID] = message.UnreadMentionIds
	}

	for _, message := range messagesSent {
		sentMap[message.ChannelID] = message.ID
		if readMap[message.ChannelID] == "" {
			readMap[message.ChannelID] = message.ID
		}
	}

	return &messageStateMaps{
		readMap:     readMap,
		sentMap:     sentMap,
		mentionsMap: mentionsMap,
	}, nil
}

func (s *userService) Setup(ctx *gin.Context) (*types.Setup, *types.APIError) {
	u, exists := ctx.Get("user")
	if !exists {
		return nil, types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	user := u.(*db.User)
	user.Password = ""

	emojis, err := s.db.GetEmojis(ctx, user.ID)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_GET_EMOJIS", "Failed to get user's emojis.", err)
	}

	friendsData, friendChannelIDs, err := s.fetchFriendsData(ctx, user.ID)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_FETCH_FRIENDS", "Failed to fetch friends data.", err)
	}

	serversData, serverChannelIDs, err := s.fetchServersData(ctx, user.ID)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_FETCH_SERVERS", "Failed to fetch servers data.", err)
	}

	allChannelIDs := append(friendChannelIDs, serverChannelIDs...)
	messageStates, err := s.buildMessageStateMaps(ctx, user.ID, allChannelIDs)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_GET_MESSAGE_STATES", "Failed to get message states.", err)
	}

	friends := s.processFriendsWithData(friendsData, messageStates)
	serversMap := s.processServersWithData(serversData, messageStates)

	return &types.Setup{
		User:    user,
		Emojis:  emojis,
		Friends: friends,
		Servers: serversMap,
	}, nil
}

func (s *userService) fetchFriendsData(ctx *gin.Context, userID string) ([]db.GetFriendsRow, []string, error) {
	dbFriends, err := s.db.GetFriends(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	channelIDs := make([]string, 0, len(dbFriends))
	for _, friend := range dbFriends {
		if friend.ChannelID.String != "" {
			channelIDs = append(channelIDs, friend.ChannelID.String)
		}
	}

	return dbFriends, channelIDs, nil
}

func (s *userService) processFriendsWithData(dbFriends []db.GetFriendsRow, messageStates *messageStateMaps) []types.Friend {
	friends := make([]types.Friend, 0, len(dbFriends))
	for _, friend := range dbFriends {
		friends = append(friends, types.Friend{
			GetFriendsRow:   friend,
			LastMessageRead: messageStates.readMap[friend.ChannelID.String],
			LastMessageSent: messageStates.sentMap[friend.ChannelID.String],
		})
	}
	return friends
}

type serverDataBundle struct {
	servers    []db.GetServersFromUserRow
	categories []db.ChannelCategory
	channels   []db.Channel
	roles      []db.GetRolesFromServersRow
}

func (s *userService) fetchServersData(ctx *gin.Context, userID string) (*serverDataBundle, []string, error) {
	servers, err := s.db.GetUserServers(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	if len(servers) == 0 {
		return &serverDataBundle{servers: servers}, []string{}, nil
	}
	serverIDs := make([]string, 0, len(servers))
	for _, server := range servers {
		serverIDs = append(serverIDs, server.ID)
	}

	type fetchResult struct {
		categories []db.ChannelCategory
		channels   []db.Channel
		roles      []db.GetRolesFromServersRow
		err        error
	}

	resultChan := make(chan fetchResult, 3)

	go func() {
		categories, err := s.db.GetCategoriesFromServers(ctx, serverIDs)
		resultChan <- fetchResult{categories: categories, err: err}
	}()

	go func() {
		channels, err := s.db.GetChannelsFromServers(ctx, serverIDs)
		resultChan <- fetchResult{channels: channels, err: err}
	}()

	go func() {
		roles, err := s.db.GetRolesFromServers(ctx, serverIDs)
		resultChan <- fetchResult{roles: roles, err: err}
	}()

	var allCategories []db.ChannelCategory
	var allChannels []db.Channel
	var allRoles []db.GetRolesFromServersRow

	for i := 0; i < 3; i++ {
		data := <-resultChan
		if data.err != nil {
			return nil, nil, data.err
		}
		if data.categories != nil {
			allCategories = data.categories
		}
		if data.channels != nil {
			allChannels = data.channels
		}
		if data.roles != nil {
			allRoles = data.roles
		}
	}

	channelIDs := make([]string, 0, len(allChannels))
	for _, channel := range allChannels {
		channelIDs = append(channelIDs, channel.ID)
	}

	return &serverDataBundle{
		servers:    servers,
		categories: allCategories,
		channels:   allChannels,
		roles:      allRoles,
	}, channelIDs, nil
}

func (s *userService) processServersWithData(data *serverDataBundle, messageStates *messageStateMaps) map[string]types.ServerWithCategories {
	if len(data.servers) == 0 {
		return make(map[string]types.ServerWithCategories)
	}

	categoriesByServer := utils.GroupBy(data.categories, func(c db.ChannelCategory) string { return c.ServerID })
	channelsByCategory := utils.GroupBy(data.channels, func(c db.Channel) string { return c.CategoryID.String })
	rolesByServer := utils.GroupBy(data.roles, func(r db.GetRolesFromServersRow) string { return r.ServerID })

	result := make(map[string]types.ServerWithCategories, len(data.servers))

	for _, server := range data.servers {
		categoryMap := make(map[string]types.CategoryWithChannels)

		for _, category := range categoriesByServer[server.ID] {
			channelMap := make(map[string]types.ServerChannel)

			for _, channel := range channelsByCategory[category.ID] {
				channelMap[channel.ID] = types.ServerChannel{
					Channel:         channel,
					LastMessageRead: messageStates.readMap[channel.ID],
					LastMessageSent: messageStates.sentMap[channel.ID],
					LastMentions:    messageStates.mentionsMap[channel.ID],
				}
			}

			categoryMap[category.ID] = types.CategoryWithChannels{
				ChannelCategory: category,
				Channels:        channelMap,
			}
		}

		result[server.ID] = types.ServerWithCategories{
			GetServersFromUserRow: server,
			Categories:            categoryMap,
			Roles:                 rolesByServer[server.ID],
		}
	}

	return result
}

func (s *userService) GetUserProfile(ctx *gin.Context) (*db.GetUserProfileRow, *types.APIError) {
	userID := ctx.Param("user_id")

	user, err := s.db.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_GET_USER_PROFILE", "Failed to get user profile.", err)
	}

	return &user, nil
}

func (s *userService) UploadEmojis(ctx *gin.Context, emojis []*multipart.FileHeader, shortcodes []string, body *types.UploadEmojiParams) (*[]types.EmojiResponse, *types.APIError) {
	u, exists := ctx.Get("user")
	if !exists {
		return nil, types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
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
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_UPLOAD_EMOJIS", "Failed to upload emojis.", err)
	}

	return &emojisResponse, nil
}

func (s *userService) UpdateEmoji(ctx *gin.Context, body *types.UpdateEmojiParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}

	userID := u.(*db.User).ID
	emojiID := ctx.Param("emoji_id")

	if err := s.db.UpdateEmoji(ctx, emojiID, userID, body); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_UPDATE_EMOJI", "Failed to update emoji.", err)
	}

	return nil
}

func (s *userService) DeleteEmoji(ctx *gin.Context) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}

	userID := u.(*db.User).ID
	emojiID := ctx.Param("emoji_id")

	if err := s.db.DeleteEmoji(ctx, emojiID, userID); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_DELETE_EMOJI", "Failed to delete emoji.", err)
	}

	return nil
}

func (s *userService) DeleteAccount(ctx *gin.Context) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	userID := u.(*db.User).ID

	token, err := ctx.Cookie("token")
	if err != nil {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_MISSING_TOKEN", "Session token not found.", err)
	}

	servers, err := s.db.GetUserServerIDs(ctx, userID)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_DELETE_ACCOUNT", "Failed to delete account.", err)
	}

	err = s.db.DeleteAccount(ctx, userID)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_DELETE_ACCOUNT", "Failed to delete account.", err)
	}

	err = s.broker.RemoveCachedUser(ctx, token)
	if err != nil {
		fmt.Println(types.NewAPIError(http.StatusInternalServerError, "ERR_DELETE_CACHE_USER", "Failed to delete cached user", err))
	}

	s.actors.NotifyAccountDeletion(userID, servers)

	return nil
}

func (s *userService) Sync(ctx *gin.Context, body *types.SyncParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	userID := u.(*db.User).ID

	if err := s.db.Sync(ctx, userID, body); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_SYNC", "Failed to sync.", err)
	}

	return nil
}
