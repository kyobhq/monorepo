package permissions

import (
	"backend/internal/database"
	"backend/internal/types"
	"context"
)

type Service interface {
	CheckPermission(ctx context.Context, serverID, userID string, ability types.Ability) bool
}

type service struct {
	db database.Service
}

func New(databaseService database.Service) Service {
	return &service{
		db: databaseService,
	}
}

func (s *service) CheckPermission(ctx context.Context, serverID, userID string, ability types.Ability) bool {
	ok, err := s.db.CheckPermission(ctx, serverID, userID, ability)
	if err != nil {
		return false
	}

	return ok
}
