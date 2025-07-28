package domains

import "backend/internal/database"

type ServerService interface{}

type serverService struct {
	db database.Service
}

func NewServerService(db database.Service) *serverService {
	return &serverService{
		db: db,
	}
}
