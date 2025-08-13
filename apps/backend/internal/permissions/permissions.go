package permissions

import (
	db "backend/db/gen_queries"
	"backend/internal/broker"
	"backend/internal/database"
	"backend/internal/types"
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CheckPermission(ctx *gin.Context, serverID string, ability types.Ability, ids ...string) bool
}

type service struct {
	db     database.Service
	broker broker.Service
}

func New(databaseService database.Service, brokerService broker.Service) Service {
	return &service{
		db:     databaseService,
		broker: brokerService,
	}
}

func (s *service) CheckPermission(ctx *gin.Context, serverID string, ability types.Ability, ids ...string) bool {
	user, exists := ctx.Get("user")
	if !exists {
		return false
	}
	userID := user.(*db.User).ID

	if ability == types.ManageMessages {
		authorID, err := s.db.GetMessageAuthor(ctx, ids[0])
		if err != nil {
			slog.Error("failed to get message author", "error", err)
			return false
		}

		if authorID == userID && authorID == ids[1] {
			return true
		}
	}

	abilities := s.getAbilities(ctx, serverID, userID)
	return slices.Contains(abilities, string(ability)) || slices.Contains(abilities, "OWNER") || slices.Contains(abilities, "ADMINISTRATOR")
}

func (s *service) getAbilities(ctx context.Context, serverID, userID string) []string {
	if abilities, err := s.broker.GetServerAbilities(ctx, serverID, userID); err == nil {
		return strings.Split(abilities, ",")
	}

	dbAbilities, err := s.db.GetServerAbilities(ctx, serverID, userID)
	if err != nil {
		slog.Error("failed to get server abilities", "error", err)
		return nil
	}

	s.broker.CacheServerAbilities(ctx, serverID, userID, dbAbilities)
	return dbAbilities
}
