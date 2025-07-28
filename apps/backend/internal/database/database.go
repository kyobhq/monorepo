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

	CreateUser(ctx context.Context, user *types.SignUpParams) (db.User, error)
	GetUser(ctx context.Context, input string) (db.User, error)
	SaveAuthToken(ctx context.Context, token, userID string) (db.Token, error)
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
	mainColor := pgtype.Text{String: "12,14,14", Valid: true}
	return s.queries.CreateUser(ctx, db.CreateUserParams{
		ID:          cuid2.Generate(),
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Password:    user.Password,
		Avatar:      avatarURL,
		MainColor:   mainColor,
	})
}

func (s *service) GetUser(ctx context.Context, input string) (db.User, error) {
	return s.queries.GetUser(ctx, db.GetUserParams{
		Email:    input,
		Username: input,
	})
}

func (s *service) SaveAuthToken(ctx context.Context, token, userID string) (db.Token, error) {
	return s.queries.CreateToken(ctx, db.CreateTokenParams{
		ID:       cuid2.Generate(),
		UserID:   userID,
		Token:    token,
		ExpireAt: time.Now().Add(30 * (24 * time.Hour)),
		Type:     "REMEMBER_ME_TOKEN",
	})
}

func (s *service) VerifyAuthToken(token string) (db.User, error) {
	return s.queries.VerifyToken(context.TODO(), token)
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
