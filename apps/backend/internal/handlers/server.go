package handlers

import (
	"backend/internal/domains"
	"backend/internal/types"
	"backend/internal/validation"
	"encoding/json"
	"net/http"
	"strconv"

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
	var body types.CreateServerParams

	maxFormSize := int64(1 << 20)
	if err := c.Request.ParseMultipartForm(maxFormSize); err != nil {
		types.NewAPIError(http.StatusBadRequest, "ERR_PARSE_FORM", "Failed to parse form", err).Respond(c)
		return
	}

	serverAvatar := c.Request.MultipartForm.File["avatar"]
	if err := validation.ValidateFiles(serverAvatar, validation.DefaultFileConfig); err != nil {
		err.Respond(c)
		return
	}

	body.Name = c.Request.FormValue("name")
	public, perr := strconv.ParseBool(c.Request.FormValue("public"))
	if perr != nil {
		types.NewAPIError(http.StatusBadRequest, "ERR_PARSE_NAME", "Failed to parse name from form.", perr).Respond(c)
		return
	}
	body.Public = public

	cropJSON := c.Request.FormValue("crop")
	if err := json.Unmarshal([]byte(cropJSON), &body.Crop); err != nil {
		types.NewAPIError(http.StatusBadRequest, "ERR_PARSE_CROP", "Failed to parse cropping informations from form.", perr).Respond(c)
		return
	}

	descriptionJSON := c.Request.FormValue("description")
	if descriptionJSON != "" {
		if err := json.Unmarshal([]byte(descriptionJSON), &body.Description); err != nil {
			types.NewAPIError(http.StatusBadRequest, "ERR_PARSE_DESCRIPTION", "Failed to parse description from form.", perr).Respond(c)
			return
		}
	}

	if verr := validation.Validate(&body); verr != nil {
		verr.Respond(c)
		return
	}

	server, err := h.domain.CreateServer(c, serverAvatar, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, server)
}

func (h *serverHandler) GetInformations() {
}

func (h *serverHandler) JoinServer(c *gin.Context) {
	var body types.JoinServerParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	serverID := c.Param("server_id")

	server, err := h.domain.JoinServer(c, serverID)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, server)
}

func (h *serverHandler) LeaveServer(c *gin.Context) {
	serverID := c.Param("server_id")

	err := h.domain.LeaveServer(c, serverID)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *serverHandler) CreateInvite(c *gin.Context) {
	serverID := c.Param("server_id")

	invite, err := h.domain.CreateInvite(c, serverID)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, invite)
}

func (h *serverHandler) DeleteInvite(c *gin.Context) {
	inviteID := c.Param("invite_id")

	err := h.domain.DeleteInvite(c, inviteID)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *serverHandler) EditProfile(c *gin.Context) {
	serverID := c.Param("server_id")

	var body types.UpdateServerProfileParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	err := h.domain.EditProfile(c, serverID, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *serverHandler) EditAvatar(c *gin.Context) {}

func (h *serverHandler) DeleteServer(c *gin.Context) {
	serverID := c.Param("server_id")

	err := h.domain.DeleteServer(c, serverID)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
