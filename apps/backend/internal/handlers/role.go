package handlers

import (
	"backend/internal/domains"
	"backend/internal/types"
	"backend/internal/validation"
	"net/http"

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

func (h *roleHandler) CreateOrEditRole(c *gin.Context) {
	var body types.CreateRoleParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	role, derr := h.domain.CreateOrEditRole(c, &body)
	if derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *roleHandler) DeleteRole(c *gin.Context) {
	var body types.DeleteRoleParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.DeleteRole(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *roleHandler) MoveRole(c *gin.Context) {
	var body types.MoveRoleMemberParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.MoveRole(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *roleHandler) AddRoleMember(c *gin.Context) {
	var body types.ChangeRoleMemberParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.AddRoleMember(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *roleHandler) RemoveRoleMember(c *gin.Context) {
	var body types.ChangeRoleMemberParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.RemoveRoleMember(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
