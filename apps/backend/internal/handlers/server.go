package handlers

import (
	"backend/internal/domains"

	"github.com/gin-gonic/gin"
)

type serverHandler struct {
	domain domains.ServerService
}

func NewServerHandlers(serverService domains.ServerService) *serverHandler {
	return &serverHandler{
		domain: serverService,
	}
}

func (h *serverHandler) CreateServer(c *gin.Context) {
}

func (h *serverHandler) JoinServer(c *gin.Context) {
}

func (h *serverHandler) LeaveServer(c *gin.Context) {
}

func (h *serverHandler) CreateInvite(c *gin.Context) {
}

func (h *serverHandler) DeleteInvite(c *gin.Context) {
}

func (h *serverHandler) EditProfile(c *gin.Context) {
}

func (h *serverHandler) DeleteServer(c *gin.Context) {
}
