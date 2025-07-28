package domains

import (
	"backend/internal/database"
)

type ChannelService interface{}

type channelService struct {
	db database.Service
}

func NewChannelService(db database.Service) *channelService {
	return &channelService{
		db: db,
	}
}
