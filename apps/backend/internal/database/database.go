package database

import (
	db "backend/db/gen_queries"
	"backend/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	random "math/rand/v2"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nrednav/cuid2"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close()

	GetUser(ctx context.Context, input string) (db.User, error)
	GetUserByID(ctx context.Context, userID string) (db.User, error)
	GetUserProfile(ctx context.Context, userID string) (db.GetUserProfileRow, error)
	CreateUser(ctx context.Context, user *types.SignUpParams) (db.User, error)
	UpdateUserAvatarNBanner(ctx context.Context, userID string, avatarURL, bannerURL *string) (db.User, error)
	UpdateUserEmail(ctx context.Context, userID string, body *types.UpdateEmailParams) (db.User, error)
	UpdateUserPassword(ctx context.Context, userID string, hashedPassword string) error
	UpdateUserProfile(ctx context.Context, userID string, body *types.UpdateProfileParams) (db.User, error)
	GetUserServers(ctx context.Context, userID string) ([]db.GetServersFromUserRow, error)
	GetUserServerIDs(ctx context.Context, userID string) ([]string, error)
	GetServersIDFromUser(ctx context.Context, userID string) ([]string, error)
	CreateServer(ctx context.Context, ownerID string, body *types.CreateServerParams, avatarURL *string) (*db.Server, error)
	CheckInvite(ctx context.Context, inviteCode string) (string, error)
	CreateInvite(ctx context.Context, userID, serverID string) (string, error)
	JoinServer(ctx context.Context, serverID string, userID string, position int) (*db.JoinServerRow, []db.ChannelCategory, []db.Channel, []db.GetRolesFromServerRow, []db.GetLatestMessagesSentRow, error)
	GetServer(ctx context.Context, serverID string) (db.Server, error)
	UpdateServerAvatarNBanner(ctx context.Context, serverID string, avatar, bannerURL *string) error
	UpdateServerProfile(ctx context.Context, serverID string, body *types.UpdateServerProfileParams) error
	LeaveServer(ctx context.Context, serverID string, userID string) error
	DeleteServer(ctx context.Context, userID, serverID string) error
	GetChannelsFromServers(ctx context.Context, serverIDs []string) ([]db.Channel, error)
	GetChannelsFromServer(ctx context.Context, serverID string) ([]db.Channel, error)
	GetCategoriesFromServers(ctx context.Context, serverIDs []string) ([]db.ChannelCategory, error)
	GetCategoriesFromServer(ctx context.Context, serverID string) ([]db.ChannelCategory, error)
	GetRolesFromServers(ctx context.Context, serverIDs []string) ([]db.GetRolesFromServersRow, error)
	GetRolesFromServer(ctx context.Context, serverID string) ([]db.GetRolesFromServerRow, error)
	GetUserRolesFromServers(ctx context.Context, userID string, serverIDs []string) ([]db.GetUserRolesFromServersRow, error)
	CreateCategory(ctx context.Context, body *types.CreateCategoryParams) (db.ChannelCategory, error)
	PinChannel(ctx context.Context, channelID, userID string, body *types.PinChannelParams) error
	CreateChannel(ctx context.Context, body *types.CreateChannelParams) (db.Channel, error)
	DeleteChannel(ctx context.Context, channelID string) error
	DeleteCategory(ctx context.Context, categoryID string) error
	UpsertRole(ctx context.Context, body *types.CreateRoleParams) (db.Role, error)
	DeleteRole(ctx context.Context, body *types.DeleteRoleParams) error
	AddRoleMember(ctx context.Context, body *types.ChangeRoleMemberParams) error
	RemoveRoleMember(ctx context.Context, body *types.ChangeRoleMemberParams) error
	MoveRole(ctx context.Context, body *types.MoveRoleMemberParams) error
	CheckPermission(ctx context.Context, serverID, userID string, ability types.Ability) (bool, error)
	GetServerAbilities(ctx context.Context, serverID, userID string) ([]string, error)
	UpdateChannelInformations(ctx context.Context, channelID string, body *types.EditChannelParams) error
	UpdateCategoryInformations(ctx context.Context, categoryID string, body *types.EditCategoryParams) error
	CreateMessage(ctx context.Context, userID string, body *types.CreateMessageParams) (db.Message, error)
	GetServers(ctx context.Context) ([]string, error)
	GetChannels(ctx context.Context) ([]db.GetChannelsIDsRow, error)
	GetServerInformations(ctx context.Context, userID, serverID string, userIDs []string) (db.GetServerInformationsRow, error)
	GetServerMembers(ctx context.Context, serverID string, offset int32, userIDs []string) ([]db.GetServerMembersRow, error)
	GetMessages(ctx context.Context, serverID, channelID, beforeMessageID, afterMessageID string, userIDs []string) ([]db.GetMessagesFromChannelRow, error)
	DeleteMessage(ctx context.Context, messageID string, userID string) error
	GetMessageAuthor(ctx context.Context, messageID string) (string, error)
	EditMessage(ctx context.Context, messageID string, body *types.EditMessageParams) error
	GetUserLinks(ctx context.Context, userID string) ([]json.RawMessage, error)
	GetUserFacts(ctx context.Context, userID string) ([]json.RawMessage, error)
	GetUserPassword(ctx context.Context, userID string) (string, error)
	UploadEmojis(ctx context.Context, userID string, emojis []db.CreateEmojiParams) error
	GetEmojis(ctx context.Context, userID string) ([]db.GetEmojisRow, error)
	UpdateEmoji(ctx context.Context, emojiID string, userID string, body *types.UpdateEmojiParams) error
	DeleteEmoji(ctx context.Context, emojiID string, userID string) error
	CreateFriendRequest(ctx context.Context, senderID, receiverID string) (db.Friend, error)
	AcceptFriendRequest(ctx context.Context, friendshipID, senderID, receiverID string) (*string, error)
	RemoveFriend(ctx context.Context, friendshipID, userID string) error
	GetFriends(ctx context.Context, userID string) ([]db.GetFriendsRow, error)
	GetFriendIDs(ctx context.Context, userID string) ([]string, error)
	DeleteAccount(ctx context.Context, userID string) error
	BanUser(ctx context.Context, serverID string, body *types.BanUserParams) error
	KickUser(ctx context.Context, serverID string, body *types.KickUserParams) error
	UnbanUser(ctx context.Context, serverID, userID string) error
	CheckBan(ctx context.Context, serverID, userID string) (pgtype.Text, error)
	GetBannedMembers(ctx context.Context, serverID string) ([]db.GetBannedMembersRow, error)
	SearchServerMembers(ctx context.Context, serverID, query string) ([]db.SearchServerMembersRow, error)
	Sync(ctx context.Context, userID string, body *types.SyncParams) error
	GetLatestMessagesRead(ctx context.Context, userID string) ([]db.GetLatestMessagesReadRow, error)
	GetLatestMessagesSent(ctx context.Context, channelIDs []string) ([]db.GetLatestMessagesSentRow, error)
}

type service struct {
	db      *pgxpool.Pool
	queries *db.Queries
}

var (
	database   = os.Getenv("PSQL_DB_DATABASE")
	password   = os.Getenv("PSQL_DB_PASSWORD")
	username   = os.Getenv("PSQL_DB_USERNAME")
	port       = os.Getenv("PSQL_DB_PORT")
	host       = os.Getenv("PSQL_DB_HOST")
	schema     = os.Getenv("PSQL_DB_SCHEMA")
	dbInstance *service
)

var defaultAvatarImages = []string{
	"https://d2vxg81co3irvv.cloudfront.net/h5l2kgkgl5388cn02q4qfxsp.webp",
	"https://d2vxg81co3irvv.cloudfront.net/avlu2a8o91gpvalk6gmgjx5r-avatar-dmplh84i8bft0ww7g60qrt15.webp",
	"https://d2vxg81co3irvv.cloudfront.net/wn3yuwwghg7efv5an78bojyc.webp",
	"https://d2vxg81co3irvv.cloudfront.net/ru3piutctxqesrv1xom8d8nn.webp",
	"https://d2vxg81co3irvv.cloudfront.net/rx38ak222ydoxb5kevww01ce-avatar-nf3531et6izcz2k6uyk5nm84.webp",
	"https://d2vxg81co3irvv.cloudfront.net/qilvsvq39rbmnybr4s24e86d-avatar-xtzdeamrkmtmh11odhtb9p2y.webp",
	"https://d2vxg81co3irvv.cloudfront.net/vlxth3gbqa3bmcsysy581pzy-avatar-frh262c7w15ntm8jia4jjrud.webp",
	"https://d2vxg81co3irvv.cloudfront.net/jnclk7c4el73kaqc726qlw7c-avatar-u0090mlv8r3zu5d850unedvu.webp",
	"https://d2vxg81co3irvv.cloudfront.net/u8b6rr2uz0xn3ubx41o6o6wx.webp",
	"https://d2vxg81co3irvv.cloudfront.net/f6983q524yum1ntx468j5gzu-avatar-mpkv59uvddg7jjf26gbigoas.webp",
	"https://d2vxg81co3irvv.cloudfront.net/qvn4bem4ms537cnn3kjyyv0l-avatar-rmk4i0evo7ls9v75jllzeuw0.webp",
}

func randomDefaultAvatar() string {
	idx := random.IntN(len(defaultAvatarImages))
	return defaultAvatarImages[idx]
}

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db:      conn,
		queries: db.New(conn),
	}

	return dbInstance
}

func (s *service) CreateUser(ctx context.Context, user *types.SignUpParams) (db.User, error) {
	avatarURL := pgtype.Text{String: randomDefaultAvatar(), Valid: true}
	return s.queries.CreateUser(ctx, db.CreateUserParams{
		ID:          cuid2.Generate(),
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Password:    user.Password,
		Avatar:      avatarURL,
	})
}

func (s *service) GetUser(ctx context.Context, input string) (db.User, error) {
	return s.queries.GetUser(ctx, db.GetUserParams{
		Email:    input,
		Username: input,
	})
}

func (s *service) GetUserByID(ctx context.Context, userID string) (db.User, error) {
	return s.queries.GetUserById(ctx, userID)
}

func (s *service) UpdateUserAvatarNBanner(ctx context.Context, userID string, avatarURL, bannerURL *string) (db.User, error) {
	var avatar pgtype.Text
	var banner pgtype.Text

	if avatarURL != nil {
		avatar = pgtype.Text{String: *avatarURL, Valid: true}
	} else {
		avatar = pgtype.Text{Valid: false}
	}

	if bannerURL != nil {
		banner = pgtype.Text{String: *bannerURL, Valid: true}
	} else {
		banner = pgtype.Text{Valid: false}
	}

	return s.queries.UpdateUserAvatarNBanner(ctx, db.UpdateUserAvatarNBannerParams{
		ID:     userID,
		Avatar: avatar,
		Banner: banner,
	})
}

func (s *service) GetUserServers(ctx context.Context, userID string) ([]db.GetServersFromUserRow, error) {
	return s.queries.GetServersFromUser(ctx, userID)
}

func (s *service) GetUserServerIDs(ctx context.Context, userID string) ([]string, error) {
	return s.queries.GetServerIDsFromUser(ctx, userID)
}

func (s *service) UpdateUserEmail(ctx context.Context, userID string, body *types.UpdateEmailParams) (db.User, error) {
	return s.queries.UpdateUserEmail(ctx, db.UpdateUserEmailParams{
		ID:    userID,
		Email: body.Email,
	})
}

func (s *service) UpdateUserPassword(ctx context.Context, userID string, hashedPassword string) error {
	return s.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:       userID,
		Password: hashedPassword,
	})
}

func (s *service) UpdateUserProfile(ctx context.Context, userID string, body *types.UpdateProfileParams) (db.User, error) {
	return s.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:          userID,
		Username:    body.Username,
		DisplayName: body.DisplayName,
		AboutMe:     body.About,
		Facts:       body.Facts,
		Links:       body.Links,
	})
}

func (s *service) CreateServer(ctx context.Context, ownerID string, body *types.CreateServerParams, avatarURL *string) (*db.Server, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	server, err := qtx.CreateServer(ctx, db.CreateServerParams{
		ID:          cuid2.Generate(),
		Name:        body.Name,
		Description: body.Description,
		Avatar:      pgtype.Text{String: *avatarURL, Valid: true},
		MainColor:   pgtype.Text{String: "#121214", Valid: true},
		OwnerID:     ownerID,
		Public:      body.Public,
	})
	if err != nil {
		return nil, err
	}

	role, err := qtx.UpsertRole(ctx, db.UpsertRoleParams{
		ID:        cuid2.Generate(),
		Position:  999,
		ServerID:  server.ID,
		Name:      "Default Permissions",
		Color:     "#ffffff",
		Abilities: []string{"VIEW_CHANNELS", "CREATE_INVITE", "CHANGE_NICKNAME", "SEND_MESSAGES", "ATTACH_FILES", "ADD_REACTIONS", "USE_PERSONAL_EMOJIS", "CONNECT", "SPEAK", "VIDEO"},
	})
	if err != nil {
		return nil, err
	}

	_, err = qtx.JoinServer(ctx, db.JoinServerParams{
		ID:       cuid2.Generate(),
		UserID:   ownerID,
		ServerID: server.ID,
		Position: int32(body.Position),
		Roles:    []string{role.ID},
	})
	if err != nil {
		return nil, err
	}

	return &server, tx.Commit(ctx)
}

func (s *service) GetServersIDFromUser(ctx context.Context, userID string) ([]string, error) {
	return s.queries.GetServersIDFromUser(ctx, userID)
}

func (s *service) CheckInvite(ctx context.Context, inviteCode string) (string, error) {
	return s.queries.CheckInvite(ctx, inviteCode)
}

func (s *service) CreateInvite(ctx context.Context, userID, serverID string) (string, error) {
	generateInviteCode, err := cuid2.Init(cuid2.WithLength(8))
	if err != nil {
		return "", err
	}

	invite, err := s.queries.GetOrCreateInvite(ctx, db.GetOrCreateInviteParams{
		ID:        cuid2.Generate(),
		CreatorID: userID,
		ServerID:  serverID,
		InviteID:  generateInviteCode(),
		ExpireAt:  time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return "", err
	}

	return invite.(string), nil
}

func (s *service) JoinServer(ctx context.Context, serverID string, userID string, position int) (*db.JoinServerRow, []db.ChannelCategory, []db.Channel, []db.GetRolesFromServerRow, []db.GetLatestMessagesSentRow, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	roleID, err := qtx.GetDefaultRoleID(ctx, serverID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	join, err := qtx.JoinServer(ctx, db.JoinServerParams{
		ID:       cuid2.Generate(),
		UserID:   userID,
		ServerID: serverID,
		Position: int32(position),
		Roles:    []string{roleID},
	})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	categories, err := qtx.GetCategoriesFromServer(ctx, serverID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	channels, err := qtx.GetChannelsFromServer(ctx, serverID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	roles, err := qtx.GetRolesFromServer(ctx, serverID)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	channelIDs := make([]string, len(channels))
	for i, channel := range channels {
		channelIDs[i] = channel.ID
	}

	latestMessagesSent, err := qtx.GetLatestMessagesSent(ctx, channelIDs)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return &join, categories, channels, roles, latestMessagesSent, tx.Commit(ctx)
}

func (s *service) UpdateServerAvatarNBanner(ctx context.Context, serverID string, avatarURL, bannerURL *string) error {
	var avatar pgtype.Text
	var banner pgtype.Text

	if avatarURL != nil {
		avatar = pgtype.Text{String: *avatarURL, Valid: true}
	} else {
		avatar = pgtype.Text{Valid: false}
	}

	if bannerURL != nil {
		banner = pgtype.Text{String: *bannerURL, Valid: true}
	} else {
		banner = pgtype.Text{Valid: false}
	}

	return s.queries.UpdateServerAvatarNBanner(ctx, db.UpdateServerAvatarNBannerParams{
		ID:     serverID,
		Avatar: avatar,
		Banner: banner,
	})
}

func (s *service) UpdateServerProfile(ctx context.Context, serverID string, body *types.UpdateServerProfileParams) error {
	return s.queries.UpdateServerProfile(ctx, db.UpdateServerProfileParams{
		ID:          serverID,
		Name:        body.Name,
		Description: body.Description,
		Public:      body.Public,
	})
}

func (s *service) GetServer(ctx context.Context, serverID string) (db.Server, error) {
	return s.queries.GetServer(ctx, serverID)
}

func (s *service) LeaveServer(ctx context.Context, serverID string, userID string) error {
	return s.queries.LeaveServer(ctx, db.LeaveServerParams{
		ServerID: serverID,
		UserID:   userID,
	})
}

func (s *service) DeleteServer(ctx context.Context, userID, serverID string) error {
	return s.queries.DeleteServer(ctx, db.DeleteServerParams{
		ID:      serverID,
		OwnerID: userID,
	})
}

func (s *service) GetChannelsFromServers(ctx context.Context, serverIDs []string) ([]db.Channel, error) {
	return s.queries.GetChannelsFromServers(ctx, serverIDs)
}

func (s *service) GetCategoriesFromServers(ctx context.Context, serverIDs []string) ([]db.ChannelCategory, error) {
	return s.queries.GetCategoriesFromServers(ctx, serverIDs)
}

func (s *service) GetRolesFromServers(ctx context.Context, serverIDs []string) ([]db.GetRolesFromServersRow, error) {
	return s.queries.GetRolesFromServers(ctx, serverIDs)
}

func (s *service) CreateCategory(ctx context.Context, body *types.CreateCategoryParams) (db.ChannelCategory, error) {
	return s.queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:       cuid2.Generate(),
		Position: int32(body.Position),
		ServerID: body.ServerID,
		Name:     body.Name,
		Users:    body.Users,
		Roles:    body.Roles,
		E2ee:     body.E2EE,
	})
}

func (s *service) PinChannel(ctx context.Context, channelID, userID string, body *types.PinChannelParams) error {
	return s.queries.PinChannel(ctx, db.PinChannelParams{
		ID:        cuid2.Generate(),
		Position:  int32(body.Position),
		ServerID:  body.ServerID,
		ChannelID: channelID,
		UserID:    userID,
	})
}

func (s *service) CreateChannel(ctx context.Context, body *types.CreateChannelParams) (db.Channel, error) {
	return s.queries.CreateChannel(ctx, db.CreateChannelParams{
		ID:          cuid2.Generate(),
		Position:    int32(body.Position),
		CategoryID:  pgtype.Text{String: body.CategoryID, Valid: true},
		ServerID:    body.ServerID,
		Name:        body.Name,
		Description: pgtype.Text{String: body.Description, Valid: true},
		Type:        body.Type,
		E2ee:        body.E2EE,
		Users:       body.Users,
		Roles:       body.Roles,
	})
}

func (s *service) DeleteChannel(ctx context.Context, channelID string) error {
	return s.queries.DeleteChannel(ctx, channelID)
}

func (s *service) DeleteCategory(ctx context.Context, categoryID string) error {
	return s.queries.DeleteCategory(ctx, categoryID)
}

func (s *service) UpsertRole(ctx context.Context, body *types.CreateRoleParams) (db.Role, error) {
	return s.queries.UpsertRole(ctx, db.UpsertRoleParams{
		ID:        body.RoleID,
		Position:  int32(body.Position),
		ServerID:  body.ServerID,
		Name:      body.Name,
		Color:     body.Color,
		Abilities: body.Abilities,
	})
}

func (s *service) DeleteRole(ctx context.Context, body *types.DeleteRoleParams) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	if err = qtx.DeleteRole(ctx, body.RoleID); err != nil {
		return err
	}

	if err = qtx.RemoveRoleFromAllMembers(ctx, body.RoleID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *service) AddRoleMember(ctx context.Context, body *types.ChangeRoleMemberParams) error {
	return s.queries.GiveRole(ctx, db.GiveRoleParams{
		ServerID:    body.ServerID,
		UserID:      body.UserID,
		ArrayAppend: body.RoleID,
	})
}

func (s *service) RemoveRoleMember(ctx context.Context, body *types.ChangeRoleMemberParams) error {
	return s.queries.RemoveRoleMember(ctx, db.RemoveRoleMemberParams{
		ServerID:    body.ServerID,
		UserID:      body.UserID,
		ArrayRemove: body.RoleID,
	})
}

func (s *service) MoveRole(ctx context.Context, body *types.MoveRoleMemberParams) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	if err := qtx.MoveRole(ctx, db.MoveRoleParams{
		ID:       body.MovedRoleID,
		Position: int32(body.To),
	}); err != nil {
		return err
	}

	if err := qtx.MoveRole(ctx, db.MoveRoleParams{
		ID:       body.TargetRoleID,
		Position: int32(body.From),
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *service) CheckPermission(ctx context.Context, serverID, userID string, ability types.Ability) (bool, error) {
	ok, err := s.queries.CheckPermission(ctx, db.CheckPermissionParams{
		ID:      serverID,
		OwnerID: userID,
		Column3: string(ability),
	})

	return ok == 1, err
}

func (s *service) GetServerAbilities(ctx context.Context, serverID, userID string) ([]string, error) {
	return s.queries.GetUserAbilities(ctx, db.GetUserAbilitiesParams{
		ServerID: serverID,
		UserID:   userID,
	})
}

func (s *service) UpdateChannelInformations(ctx context.Context, channelID string, body *types.EditChannelParams) error {
	return s.queries.UpdateChannelInformations(ctx, db.UpdateChannelInformationsParams{
		ID:          channelID,
		Name:        body.Name,
		Description: pgtype.Text{String: body.Description, Valid: true},
		Users:       body.Users,
		Roles:       body.Roles,
	})
}

func (s *service) UpdateCategoryInformations(ctx context.Context, categoryID string, body *types.EditCategoryParams) error {
	return s.queries.UpdateCategoryInformations(ctx, db.UpdateCategoryInformationsParams{
		ID:    categoryID,
		Name:  body.Name,
		Users: body.Users,
		Roles: body.Roles,
	})
}

func (s *service) CreateMessage(ctx context.Context, userID string, body *types.CreateMessageParams) (db.Message, error) {
	return s.queries.CreateMessage(ctx, db.CreateMessageParams{
		ID:               cuid2.Generate(),
		AuthorID:         userID,
		ServerID:         body.ServerID,
		ChannelID:        body.ChannelID,
		Content:          body.Content,
		Everyone:         body.Everyone,
		MentionsUsers:    body.MentionsUsers,
		MentionsRoles:    body.MentionsRoles,
		MentionsChannels: body.MentionsChannels,
		Attachments:      body.Attachments,
	})
}

func (s *service) GetServers(ctx context.Context) ([]string, error) {
	return s.queries.GetServersIDs(ctx)
}

func (s *service) GetChannels(ctx context.Context) ([]db.GetChannelsIDsRow, error) {
	return s.queries.GetChannelsIDs(ctx)
}

func (s *service) GetServerInformations(ctx context.Context, userID, serverID string, userIDs []string) (db.GetServerInformationsRow, error) {
	return s.queries.GetServerInformations(ctx, db.GetServerInformationsParams{
		ID:      serverID,
		Column2: userIDs,
		UserID:  userID,
	})
}

func (s *service) GetServerMembers(ctx context.Context, serverID string, offset int32, userIDs []string) ([]db.GetServerMembersRow, error) {
	return s.queries.GetServerMembers(ctx, db.GetServerMembersParams{
		ServerID: serverID,
		Offset:   offset,
		Column3:  userIDs,
	})
}

func (s *service) GetMessages(ctx context.Context, serverID, channelID, beforeMessageID, afterMessageID string, userIDs []string) ([]db.GetMessagesFromChannelRow, error) {
	return s.queries.GetMessagesFromChannel(ctx, db.GetMessagesFromChannelParams{
		ServerID:  serverID,
		ChannelID: channelID,
		Column3:   beforeMessageID,
		Column4:   afterMessageID,
		Column5:   userIDs,
	})
}

func (s *service) DeleteMessage(ctx context.Context, messageID string, userID string) error {
	return s.queries.DeleteMessage(ctx, db.DeleteMessageParams{
		ID:       messageID,
		AuthorID: userID,
	})
}

func (s *service) GetMessageAuthor(ctx context.Context, messageID string) (string, error) {
	return s.queries.GetMessageAuthor(ctx, messageID)
}

func (s *service) EditMessage(ctx context.Context, messageID string, body *types.EditMessageParams) error {
	return s.queries.UpdateMessage(ctx, db.UpdateMessageParams{
		ID:               messageID,
		Content:          body.Content,
		Everyone:         body.Everyone,
		MentionsUsers:    body.MentionsUsers,
		MentionsChannels: body.MentionsChannels,
	})
}

func (s *service) GetUserLinks(ctx context.Context, userID string) ([]json.RawMessage, error) {
	return s.queries.GetUserLinks(ctx, userID)
}

func (s *service) GetUserFacts(ctx context.Context, userID string) ([]json.RawMessage, error) {
	return s.queries.GetUserFacts(ctx, userID)
}

func (s *service) GetUserPassword(ctx context.Context, userID string) (string, error) {
	return s.queries.GetUserPassword(ctx, userID)
}

func (s *service) GetUserProfile(ctx context.Context, userID string) (db.GetUserProfileRow, error) {
	return s.queries.GetUserProfile(ctx, userID)
}

func (s *service) GetChannelsFromServer(ctx context.Context, serverID string) ([]db.Channel, error) {
	return s.queries.GetChannelsFromServer(ctx, serverID)
}

func (s *service) GetCategoriesFromServer(ctx context.Context, serverID string) ([]db.ChannelCategory, error) {
	return s.queries.GetCategoriesFromServer(ctx, serverID)
}

func (s *service) GetRolesFromServer(ctx context.Context, serverID string) ([]db.GetRolesFromServerRow, error) {
	return s.queries.GetRolesFromServer(ctx, serverID)
}

func (s *service) GetUserRolesFromServers(ctx context.Context, userID string, serverIDs []string) ([]db.GetUserRolesFromServersRow, error) {
	return s.queries.GetUserRolesFromServers(ctx, db.GetUserRolesFromServersParams{
		UserID:  userID,
		Column2: serverIDs,
	})
}

func (s *service) UploadEmojis(ctx context.Context, userID string, emojis []db.CreateEmojiParams) error {
	if _, err := s.queries.CreateEmoji(ctx, emojis); err != nil {
		return err
	}

	return nil
}

func (s *service) GetEmojis(ctx context.Context, userID string) ([]db.GetEmojisRow, error) {
	return s.queries.GetEmojis(ctx, userID)
}

func (s *service) DeleteEmoji(ctx context.Context, emojiID string, userID string) error {
	return s.queries.DeleteEmoji(ctx, db.DeleteEmojiParams{
		ID:     emojiID,
		UserID: userID,
	})
}

func (s *service) UpdateEmoji(ctx context.Context, emojiID string, userID string, body *types.UpdateEmojiParams) error {
	return s.queries.UpdateEmoji(ctx, db.UpdateEmojiParams{
		ID:        emojiID,
		UserID:    userID,
		Shortcode: body.Shortcode,
	})
}

func (s *service) CreateFriendRequest(ctx context.Context, senderID, receiverID string) (db.Friend, error) {
	return s.queries.AddFriend(ctx, db.AddFriendParams{
		ID:         cuid2.Generate(),
		SenderID:   senderID,
		ReceiverID: receiverID,
	})
}

func (s *service) AcceptFriendRequest(ctx context.Context, friendshipID, senderID, receiverID string) (*string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	err = qtx.AcceptFriend(ctx, db.AcceptFriendParams{
		ID:         friendshipID,
		ReceiverID: receiverID,
	})
	if err != nil {
		return nil, err
	}

	channel, err := qtx.CreateChannel(ctx, db.CreateChannelParams{
		ID:           cuid2.Generate(),
		ServerID:     "global",
		FriendshipID: pgtype.Text{String: friendshipID, Valid: true},
		Name:         "",
		Type:         "dm",
		E2ee:         false,
		Users:        []string{senderID, receiverID},
		Description:  pgtype.Text{String: "", Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &channel.ID, tx.Commit(ctx)
}

func (s *service) RemoveFriend(ctx context.Context, friendshipID, userID string) error {
	return s.queries.DeleteFriend(ctx, db.DeleteFriendParams{
		ID:         friendshipID,
		ReceiverID: userID,
	})
}

func (s *service) GetFriends(ctx context.Context, userID string) ([]db.GetFriendsRow, error) {
	return s.queries.GetFriends(ctx, userID)
}

func (s *service) GetFriendIDs(ctx context.Context, userID string) ([]string, error) {
	return s.queries.GetFriendIDs(ctx, userID)
}

func (s *service) DeleteAccount(ctx context.Context, userID string) error {
	return s.queries.DeleteUser(ctx, userID)
}

func (s *service) BanUser(ctx context.Context, serverID string, body *types.BanUserParams) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	qtx.BanUser(ctx, db.BanUserParams{
		UserID:    body.UserID,
		ServerID:  serverID,
		BanReason: pgtype.Text{String: body.Reason, Valid: true},
	})

	qtx.DeleteServerMessages(ctx, db.DeleteServerMessagesParams{
		ServerID: serverID,
		AuthorID: body.UserID,
	})

	return tx.Commit(ctx)
}

func (s *service) KickUser(ctx context.Context, serverID string, body *types.KickUserParams) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	qtx.KickUser(ctx, db.KickUserParams{
		UserID:   body.UserID,
		ServerID: serverID,
	})

	qtx.DeleteServerMessages(ctx, db.DeleteServerMessagesParams{
		ServerID: serverID,
		AuthorID: body.UserID,
	})

	return tx.Commit(ctx)
}

func (s *service) UnbanUser(ctx context.Context, serverID, userID string) error {
	return s.queries.KickUser(ctx, db.KickUserParams{
		UserID:   userID,
		ServerID: serverID,
	})
}

func (s *service) CheckBan(ctx context.Context, serverID, userID string) (pgtype.Text, error) {
	return s.queries.CheckBan(ctx, db.CheckBanParams{
		UserID:   userID,
		ServerID: serverID,
	})
}

func (s *service) GetBannedMembers(ctx context.Context, serverID string) ([]db.GetBannedMembersRow, error) {
	return s.queries.GetBannedMembers(ctx, serverID)
}

func (s *service) SearchServerMembers(ctx context.Context, serverID, query string) ([]db.SearchServerMembersRow, error) {
	return s.queries.SearchServerMembers(ctx, db.SearchServerMembersParams{
		ServerID:    serverID,
		DisplayName: "%" + query + "%",
	})
}

func (s *service) Sync(ctx context.Context, userID string, body *types.SyncParams) error {
	return s.queries.SaveUnreadMessagesState(ctx, db.SaveUnreadMessagesStateParams{
		UserID:             userID,
		ChannelIds:         body.ChannelIDs,
		LastReadMessageIds: body.LastMessageIDs,
		UnreadMentionIds:   body.MentionIDs,
	})
}

func (s *service) GetLatestMessagesRead(ctx context.Context, userID string) ([]db.GetLatestMessagesReadRow, error) {
	return s.queries.GetLatestMessagesRead(ctx, userID)
}

func (s *service) GetLatestMessagesSent(ctx context.Context, channelIDs []string) ([]db.GetLatestMessagesSentRow, error) {
	return s.queries.GetLatestMessagesSent(ctx, channelIDs)
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stat()
	stats["open_connections"] = strconv.Itoa(int(dbStats.TotalConns()))
	stats["acquire_count"] = strconv.FormatInt(dbStats.AcquireCount(), 10)
	stats["acquired_conns"] = strconv.Itoa(int(dbStats.AcquiredConns()))
	stats["idle_conns"] = strconv.Itoa(int(dbStats.IdleConns()))
	stats["max_conns"] = strconv.Itoa(int(dbStats.MaxConns()))

	// Evaluate stats to provide a health message
	if dbStats.TotalConns() > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.AcquiredConns() == dbStats.MaxConns() {
		stats["message"] = "Connection pool is at maximum capacity."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() {
	log.Printf("Disconnected from database: %s", database)
	s.db.Close()
}
