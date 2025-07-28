package domains

import "backend/internal/database"

type UserService interface{}

type userService struct {
	db database.Service
}

func NewUserService(db database.Service) *userService {
	return &userService{
		db: db,
	}
}
