package domains

import (
	"backend/internal/database"
	"backend/internal/permissions"
)

type RoleService interface{}

type roleService struct {
	db          database.Service
	permissions permissions.Service
}

func NewRoleService(db database.Service, permissions permissions.Service) *roleService {
	return &roleService{
		db:          db,
		permissions: permissions,
	}
}
