package handlers

import (
	"backend/internal/domains"

	"github.com/gin-gonic/gin"
)

type roleHandler struct {
	domain domains.RoleService
}

func NewRoleHandlers(roleService domains.RoleService) *roleHandler {
	return &roleHandler{
		domain: roleService,
	}
}

func (h *roleHandler) GetRoles(c *gin.Context) {
}

func (h *roleHandler) CreateRole(c *gin.Context) {
}

func (h *roleHandler) EditRole(c *gin.Context) {
}

func (h *roleHandler) DeleteRole(c *gin.Context) {
}

func (h *roleHandler) MoveRole(c *gin.Context) {
}

func (h *roleHandler) AddRoleMember(c *gin.Context) {
}

func (h *roleHandler) RemoveRoleMember(c *gin.Context) {
}
