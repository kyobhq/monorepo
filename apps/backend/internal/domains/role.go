package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/actors"
	"backend/internal/database"
	"backend/internal/permissions"
	"backend/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleService interface {
	CreateOrEditRole(ctx *gin.Context, body *types.CreateRoleParams) (*db.Role, *types.APIError)
	DeleteRole(ctx *gin.Context, body *types.DeleteRoleParams) *types.APIError
	AddRoleMember(ctx *gin.Context, body *types.ChangeRoleMemberParams) *types.APIError
	RemoveRoleMember(ctx *gin.Context, body *types.ChangeRoleMemberParams) *types.APIError
	MoveRole(ctx *gin.Context, body *types.MoveRoleMemberParams) *types.APIError
}

type roleService struct {
	db          database.Service
	actors      actors.Service
	permissions permissions.Service
}

func NewRoleService(db database.Service, actors actors.Service, permissions permissions.Service) *roleService {
	return &roleService{
		db:          db,
		actors:      actors,
		permissions: permissions,
	}
}

func (s *roleService) CreateOrEditRole(ctx *gin.Context, body *types.CreateRoleParams) (*db.Role, *types.APIError) {
	if allowed := s.permissions.CheckPermission(ctx, body.ServerID, types.ManageRoles); !allowed {
		return nil, &types.APIError{
			Status:  http.StatusForbidden,
			Code:    "ERR_FORBIDDEN",
			Cause:   "",
			Message: "Forbidden to create a role",
		}
	}

	role, err := s.db.UpsertRole(ctx, body)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CREATE_OR_EDIT_ROLE",
			Cause:   err.Error(),
			Message: "Failed to create or edit a role",
		}
	}

	s.actors.CreateOrEditRole(role)

	return &role, nil
}

func (s *roleService) DeleteRole(ctx *gin.Context, body *types.DeleteRoleParams) *types.APIError {
	if allowed := s.permissions.CheckPermission(ctx, body.ServerID, types.ManageRoles); !allowed {
		return &types.APIError{
			Status:  http.StatusForbidden,
			Code:    "ERR_FORBIDDEN",
			Cause:   "",
			Message: "Forbidden to create a role",
		}
	}

	if err := s.db.DeleteRole(ctx, body); err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CREATE_ROLE",
			Cause:   err.Error(),
			Message: "Failed to create a role",
		}
	}

	s.actors.RemoveRole(body)

	return nil
}

func (s *roleService) AddRoleMember(ctx *gin.Context, body *types.ChangeRoleMemberParams) *types.APIError {
	if allowed := s.permissions.CheckPermission(ctx, body.ServerID, types.ManageRoles); !allowed {
		return &types.APIError{
			Status:  http.StatusForbidden,
			Code:    "ERR_FORBIDDEN",
			Cause:   "",
			Message: "Forbidden to create a role",
		}
	}

	if err := s.db.AddRoleMember(ctx, body); err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_ADD_ROLE_TO_MEMBER",
			Cause:   err.Error(),
			Message: "Failed to add role to a member",
		}
	}

	s.actors.AddRoleMember(body)

	return nil
}

func (s *roleService) RemoveRoleMember(ctx *gin.Context, body *types.ChangeRoleMemberParams) *types.APIError {
	if allowed := s.permissions.CheckPermission(ctx, body.ServerID, types.ManageRoles); !allowed {
		return &types.APIError{
			Status:  http.StatusForbidden,
			Code:    "ERR_FORBIDDEN",
			Cause:   "",
			Message: "Forbidden to create a role",
		}
	}

	if err := s.db.RemoveRoleMember(ctx, body); err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_REMOVE_ROLE_FROM_MEMBER",
			Cause:   err.Error(),
			Message: "Failed to remove role from a member",
		}
	}

	s.actors.RemoveRoleMember(body)

	return nil
}

func (s *roleService) MoveRole(ctx *gin.Context, body *types.MoveRoleMemberParams) *types.APIError {
	if allowed := s.permissions.CheckPermission(ctx, body.ServerID, types.ManageRoles); !allowed {
		return &types.APIError{
			Status:  http.StatusForbidden,
			Code:    "ERR_FORBIDDEN",
			Cause:   "",
			Message: "Forbidden to create a role",
		}
	}

	if err := s.db.MoveRole(ctx, body); err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_MOVE_ROLE",
			Cause:   err.Error(),
			Message: "Failed to move role",
		}
	}

	s.actors.MoveRole(body)

	return nil
}
