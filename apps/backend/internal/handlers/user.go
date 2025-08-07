package handlers

import (
	"backend/internal/domains"
	"backend/internal/types"
	"backend/internal/validation"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	domain domains.UserService
}

func NewUserHandlers(userService domains.UserService) *userHandler {
	return &userHandler{
		domain: userService,
	}
}

func (h *userHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("user_id")

	user, err := h.domain.GetUserProfile(c, userID)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) Setup(c *gin.Context) {
	setup, err := h.domain.Setup(c)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, setup)
}

func (h *userHandler) UpdateEmail(c *gin.Context) {
	var body types.UpdateEmailParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	err := h.domain.UpdateEmail(c, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *userHandler) UpdatePassword(c *gin.Context) {
	var body types.UpdatePasswordParams

	if verr := validation.ParseAndValidate(c.Request, &body); verr != nil {
		verr.Respond(c)
		return
	}

	if body.New != body.Confirm {
		types.NewAPIError(http.StatusBadRequest, "ERR_PASSWORD_MISMATCH", "Passwords do not match", nil).Respond(c)
		return
	}

	err := h.domain.UpdatePassword(c, &body)
	if err != nil {
		err.Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *userHandler) UpdateProfile(c *gin.Context) {
	var body types.UpdateProfileParams

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

func (h *userHandler) UpdateAvatar(c *gin.Context) {
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
