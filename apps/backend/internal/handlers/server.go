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

func (h *serverHandler) GetInformations(c *gin.Context) {
	serverInformations, err := h.domain.GetInformations(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, serverInformations)
}

func (h *serverHandler) GetMembers(c *gin.Context) {
	members, err := h.domain.GetMembers(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *serverHandler) GetBannedMembers(c *gin.Context) {
	bans, err := h.domain.GetBannedMembers(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, bans)
}

func (h *serverHandler) JoinServer(c *gin.Context) {
	var body types.JoinServerParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	server, err := h.domain.JoinServer(c, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, server)
}

func (h *serverHandler) LeaveServer(c *gin.Context) {
	err := h.domain.LeaveServer(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *serverHandler) CreateInvite(c *gin.Context) {
	invite, err := h.domain.CreateInvite(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, invite)
}

func (h *serverHandler) DeleteInvite(c *gin.Context) {
	err := h.domain.DeleteInvite(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *serverHandler) UpdateProfile(c *gin.Context) {
	var body types.UpdateServerProfileParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	err := h.domain.UpdateProfile(c, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *serverHandler) UpdateAvatar(c *gin.Context) {
	var body types.UpdateAvatarParams

	maxFormSize := int64(1 << 20)
	if err := c.Request.ParseMultipartForm(maxFormSize); err != nil {
		types.NewAPIError(http.StatusBadRequest, "ERR_PARSE_FORM", "Failed to parse form", err).Respond(c)
		return
	}

	avatar := c.Request.MultipartForm.File["avatar"]
	if len(avatar) > 0 {
		if err := validation.ValidateFiles(avatar, validation.DefaultFileConfig); err != nil {
			err.Respond(c)
			return
		}
	}

	banner := c.Request.MultipartForm.File["banner"]
	if len(banner) > 0 {
		if err := validation.ValidateFiles(banner, validation.DefaultFileConfig); err != nil {
			err.Respond(c)
			return
		}
	}

	cropAvatarJSON := c.Request.FormValue("crop_avatar")
	if err := json.Unmarshal([]byte(cropAvatarJSON), &body.CropAvatar); err != nil {
		types.NewAPIError(http.StatusBadRequest, "ERR_PARSE_CROP", "Failed to parse cropping informations from form.", err).Respond(c)
		return
	}

	crapBannerJSON := c.Request.FormValue("crop_banner")
	if err := json.Unmarshal([]byte(crapBannerJSON), &body.CropBanner); err != nil {
		types.NewAPIError(http.StatusBadRequest, "ERR_PARSE_CROP", "Failed to parse cropping informations from form.", err).Respond(c)
		return
	}

	if verr := validation.Validate(&body); verr != nil {
		verr.Respond(c)
		return
	}

	avatarURL, bannerURL, err := h.domain.UpdateAvatar(c, avatar, banner, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar": avatarURL, "banner": bannerURL})
}

func (h *serverHandler) DeleteServer(c *gin.Context) {
	err := h.domain.DeleteServer(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *serverHandler) BanUser(c *gin.Context) {
	var body types.BanUserParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.BanUser(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *serverHandler) UnbanUser(c *gin.Context) {
	err := h.domain.UnbanUser(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *serverHandler) KickUser(c *gin.Context) {
	var body types.KickUserParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if derr := h.domain.KickUser(c, &body); derr != nil {
		derr.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
