package domains

import (
	"backend/internal/database"
)

type FriendService interface{}

type friendService struct {
	db database.Service
}

func NewFriendService(db database.Service) *friendService {
	return &friendService{
		db: db,
	}
}
