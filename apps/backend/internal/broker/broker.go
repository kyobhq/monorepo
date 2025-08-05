package broker

import (
	db "backend/db/gen_queries"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// PublishTo is used to send a message to a specific channel.
	// Channels follows this pattern: actor:channel (e.g. server:123, channel:123).
	// Messages are protobuf encoded.
	PublishTo(channel string, message []byte) error

	// SubcribeTo is used to subscribe to a specific channel.
	// Channels follows this pattern: actor:channel (e.g. server:123, channel:123).
	SubcribeTo(channels ...string) *redis.PubSub

	CacheUser(ctx context.Context, token string, user db.User) error
	GetCachedUser(ctx context.Context, token string) (*db.User, error)
	RefreshCachedUser(ctx context.Context, token string, user db.User) error
	RemoveCachedUser(ctx context.Context, token string) error

	CacheServerAbilities(ctx context.Context, serverID, userID string, abilities []string) error
	GetServerAbilities(ctx context.Context, serverID, userID string) (string, error)
}

type service struct {
	db *redis.Client
}

type PresenceInfo struct {
	PID    *actor.PID `json:"pid"`
	Status string     `json:"status"`
}

var (
	password   = os.Getenv("DRAGONFLY_DB_PASSWORD")
	port       = os.Getenv("DRAGONFLY_DB_PORT")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	options := &redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", port),
		Password: password,
		DB:       0,
	}

	db := redis.NewClient(options)
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) SubcribeTo(channels ...string) *redis.PubSub {
	return s.db.Subscribe(context.TODO(), channels...)
}

func (s *service) PublishTo(channel string, message []byte) error {
	return s.db.Publish(context.TODO(), channel, message).Err()
}

func (s *service) CacheUser(ctx context.Context, token string, user db.User) error {
	user.Password = ""

	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	s.db.Set(ctx, "user:"+token, userJSON, 30*(24*time.Hour))

	return nil
}

func (s *service) GetCachedUser(ctx context.Context, token string) (*db.User, error) {
	res := s.db.Get(ctx, "user:"+token)
	userJSON, err := res.Result()
	if err != nil {
		return nil, err
	}

	var user db.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) RefreshCachedUser(ctx context.Context, token string, user db.User) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	key := "user:" + token
	ttl, err := s.db.TTL(ctx, key).Result()
	if err != nil {
		return err
	}

	if ttl < 0 {
		ttl = 30 * (24 * time.Hour)
	}

	if err = s.db.Set(ctx, key, userJSON, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (s *service) RemoveCachedUser(ctx context.Context, token string) error {
	res := s.db.Del(ctx, "user:"+token)
	_, err := res.Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CacheServerAbilities(ctx context.Context, serverID, userID string, abilities []string) error {
	key := fmt.Sprintf("roles:%s:%s", serverID, userID)
	res := s.db.Set(ctx, key, strings.Join(abilities, ","), 10*time.Minute)
	_, err := res.Result()

	return err
}

func (s *service) GetServerAbilities(ctx context.Context, serverID, userID string) (string, error) {
	key := fmt.Sprintf("roles:%s:%s", serverID, userID)
	res := s.db.Get(ctx, key)

	return res.Result()
}

// Health checks the health of the broker connection by pinging the broker.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the broker
	pong := s.db.Ping(ctx)
	if pong.Err() != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("broker down: %v", pong.Err())
		log.Fatalf("db down: %v", pong.Err()) // Log the error and terminate the program
		return stats
	}

	// Broker is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get broker stats (like total connections, stale, idle, etc.)
	poolStats := s.db.PoolStats()
	stats["hits"] = strconv.FormatUint(uint64(poolStats.Hits), 10)
	stats["misses"] = strconv.FormatUint(uint64(poolStats.Misses), 10)
	stats["timeouts"] = strconv.FormatUint(uint64(poolStats.Timeouts), 10)
	stats["total_connections"] = strconv.FormatUint(uint64(poolStats.TotalConns), 10)
	stats["idle_connections"] = strconv.FormatUint(uint64(poolStats.IdleConns), 10)
	stats["stale_connections"] = strconv.FormatUint(uint64(poolStats.StaleConns), 10)

	// Evaluate stats to provide a health message
	if poolStats.TotalConns > 40 {
		stats["message"] = "The broker is experiencing heavy load."
	}

	if poolStats.Timeouts > 1000 {
		stats["message"] = "High timeouts indicate potential connection pool bottlenecks."
	}

	if poolStats.StaleConns > poolStats.TotalConns/2 {
		stats["message"] = "Many stale connections; consider tuning pool timeouts or network settings."
	}

	if poolStats.Hits+poolStats.Misses > 10 && poolStats.Misses > poolStats.Hits/2 {
		stats["message"] = "High miss rate—expand the connection pool or optimize queries."
	}

	info := s.db.Info(ctx, "memory").Val()
	stats["memory_info"] = info

	return stats
}

// Close closes the broker connection.
// It logs a message indicating the disconnection from the specific broker.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from broker")
	return s.db.Close()
}
