package domains

import "backend/internal/database"

type RoleService interface{}

type roleService struct {
	db database.Service
}

func NewRoleService(db database.Service) *roleService {
	return &roleService{
		db: db,
	}
}
