package database

import (
	db "backend/db/gen_queries"
	"backend/internal/types"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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
	CreateUser(ctx context.Context, user *types.SignUpParams) (db.User, error)
	UpdateUserAvatarNBanner(ctx context.Context, userID string, body *types.UpdateAvatarParams) error
	UpdateUserAccount(ctx context.Context, userID string, body *types.UpdateAccountParams) error
	UpdateUserPassword(ctx context.Context, userID string, password string) error
	UpdateUserProfile(ctx context.Context, userID string, body *types.UpdateProfileParams) error
	GetUserServers(ctx context.Context, userID string) ([]db.GetServersFromUserRow, error)
	CreateServer(ctx context.Context, ownerID string, body *types.CreateServerParams, avatarURL *string) (*db.Server, error)
	CheckInvite(ctx context.Context, inviteCode string) (string, error)
	CreateInvite(ctx context.Context, serverID string) (string, error)
	JoinServer(ctx context.Context, serverID string, userID string) error
	GetServer(ctx context.Context, serverID string) (db.Server, error)
	UpdateServerAvatarNBanner(ctx context.Context, serverID string, body *types.UpdateServerAvatarParams) error
	UpdateServerProfile(ctx context.Context, serverID string, body *types.UpdateServerProfileParams) error
	LeaveServer(ctx context.Context, serverID string, userID string) error
	DeleteServer(ctx context.Context, serverID string) (pgconn.CommandTag, error)
	GetChannelsFromServers(ctx context.Context, serverIDs []string) ([]db.Channel, error)
	GetCategoriesFromServers(ctx context.Context, serverIDs []string) ([]db.ChannelCategory, error)
	GetRolesFromServers(ctx context.Context, serverIDs []string) ([]db.GetRolesFromServersRow, error)
	CreateCategory(ctx context.Context, body *types.CreateCategoryParams) (db.ChannelCategory, error)
	PinChannel(ctx context.Context, channelID, userID string, body *types.PinChannelParams) error
	CreateChannel(ctx context.Context, body *types.CreateChannelParams) (db.Channel, error)
	DeleteChannel(ctx context.Context, channelID string) error
	DeleteCategory(ctx context.Context, categoryID string) error
	CreateRole(ctx context.Context, body *types.CreateRoleParams) (db.Role, error)
	CheckPermission(ctx context.Context, serverID, userID string, ability types.Ability) (bool, error)
	GetServerAbilities(ctx context.Context, serverID, userID string) ([]string, error)
	UpdateChannelInformations(ctx context.Context, channelID string, body *types.EditChannelParams) error
	CreateMessage(ctx context.Context, userID string, body *types.CreateMessageParams) (db.Message, error)
	GetServers(ctx context.Context) ([]string, error)
	GetChannels(ctx context.Context) ([]db.GetChannelsIDsRow, error)
	GetServerInformations(ctx context.Context, serverID string) (db.GetServerInformationsRow, error)
	GetMessages(ctx context.Context, channelID string) ([]db.GetMessagesFromChannelRow, error)
	DeleteMessage(ctx context.Context, messageID string, userID string) error
	GetMessageAuthor(ctx context.Context, messageID string) (string, error)
	EditMessage(ctx context.Context, messageID string, body *types.EditMessageParams) error
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
	avatarURL := pgtype.Text{String: "https://i.pinimg.com/1200x/ef/cf/e5/efcfe5321149cb491399bd159586a2ec.jpg", Valid: true}
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

func (s *service) UpdateUserAvatarNBanner(ctx context.Context, userID string, body *types.UpdateAvatarParams) error {
	avatarURL := pgtype.Text{String: body.Avatar, Valid: true}
	bannerURL := pgtype.Text{String: body.Banner, Valid: true}
	mainColor := pgtype.Text{String: body.MainColor, Valid: true}

	return s.queries.UpdateUserAvatarNBanner(ctx, db.UpdateUserAvatarNBannerParams{
		ID:        userID,
		Avatar:    avatarURL,
		Banner:    bannerURL,
		MainColor: mainColor,
	})
}

func (s *service) GetUserServers(ctx context.Context, userID string) ([]db.GetServersFromUserRow, error) {
	return s.queries.GetServersFromUser(ctx, userID)
}

func (s *service) UpdateUserAccount(ctx context.Context, userID string, body *types.UpdateAccountParams) error {
	return s.queries.UpdateUserInformations(ctx, db.UpdateUserInformationsParams{
		ID:       userID,
		Email:    body.Email,
		Username: body.Username,
	})
}

func (s *service) UpdateUserPassword(ctx context.Context, userID string, password string) error {
	return s.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:       userID,
		Password: password,
	})
}

func (s *service) UpdateUserProfile(ctx context.Context, userID string, body *types.UpdateProfileParams) error {
	return s.queries.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:          userID,
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
		MainColor:   pgtype.Text{String: "12,14,14", Valid: true},
		OwnerID:     ownerID,
		Public:      body.Public,
	})
	if err != nil {
		return nil, err
	}

	err = qtx.JoinServer(ctx, db.JoinServerParams{
		ID:       cuid2.Generate(),
		UserID:   ownerID,
		ServerID: server.ID,
		Position: int32(body.Position),
	})
	if err != nil {
		return nil, err
	}

	return &server, tx.Commit(ctx)
}

func (s *service) CheckInvite(ctx context.Context, inviteCode string) (string, error) {
	return s.queries.CheckInvite(ctx, inviteCode)
}

func (s *service) CreateInvite(ctx context.Context, serverID string) (string, error) {
	return s.queries.CreateInvite(ctx, db.CreateInviteParams{
		ID:       cuid2.Generate(),
		ServerID: serverID,
		InviteID: cuid2.Generate(),
		ExpireAt: time.Now().Add(time.Hour * 24),
	})
}

func (s *service) JoinServer(ctx context.Context, serverID string, userID string) error {
	return s.queries.JoinServer(ctx, db.JoinServerParams{
		ID:       serverID,
		UserID:   userID,
		ServerID: serverID,
	})
}

func (s *service) UpdateServerAvatarNBanner(ctx context.Context, serverID string, body *types.UpdateServerAvatarParams) error {
	avatarURL := pgtype.Text{String: body.Avatar, Valid: true}
	bannerURL := pgtype.Text{String: body.Banner, Valid: true}
	mainColor := pgtype.Text{String: body.MainColor, Valid: true}

	return s.queries.UpdateServerAvatarNBanner(ctx, db.UpdateServerAvatarNBannerParams{
		ID:        serverID,
		Avatar:    avatarURL,
		Banner:    bannerURL,
		MainColor: mainColor,
	})
}

func (s *service) UpdateServerProfile(ctx context.Context, serverID string, body *types.UpdateServerProfileParams) error {
	return s.queries.UpdateServerProfile(ctx, db.UpdateServerProfileParams{
		ID:          serverID,
		Name:        body.Name,
		Description: body.Description,
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

func (s *service) DeleteServer(ctx context.Context, serverID string) (pgconn.CommandTag, error) {
	return s.queries.DeleteServer(ctx, db.DeleteServerParams{
		ID: serverID,
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
		CategoryID:  body.CategoryID,
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

func (s *service) CreateRole(ctx context.Context, body *types.CreateRoleParams) (db.Role, error) {
	return s.queries.CreateRole(ctx, db.CreateRoleParams{
		ID:        cuid2.Generate(),
		Position:  int32(body.Position),
		ServerID:  body.ServerID,
		Name:      body.Name,
		Color:     body.Color,
		Abilities: body.Abilities,
	})
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

func (s *service) GetServerInformations(ctx context.Context, serverID string) (db.GetServerInformationsRow, error) {
	return s.queries.GetServerInformations(ctx, serverID)
}

func (s *service) GetMessages(ctx context.Context, channelID string) ([]db.GetMessagesFromChannelRow, error) {
	return s.queries.GetMessagesFromChannel(ctx, channelID)
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
