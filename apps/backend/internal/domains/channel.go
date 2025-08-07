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

type ChannelService interface {
	CreateCategory(c *gin.Context, body *types.CreateCategoryParams) (*db.ChannelCategory, *types.APIError)
	CreateChannel(c *gin.Context, body *types.CreateChannelParams) (*db.Channel, *types.APIError)
	PinChannel(c *gin.Context, body *types.PinChannelParams) *types.APIError
	DeleteChannel(c *gin.Context, body *types.DeleteChannelParams) *types.APIError
	DeleteCategory(c *gin.Context, body *types.DeleteCategoryParams) *types.APIError
	EditChannel(c *gin.Context, body *types.EditChannelParams) *types.APIError
}

type channelService struct {
	db          database.Service
	actors      actors.Service
	permissions permissions.Service
}

func NewChannelService(db database.Service, actors actors.Service, permissions permissions.Service) *channelService {
	return &channelService{
		db:          db,
		actors:      actors,
		permissions: permissions,
	}
}

func (s *channelService) CreateCategory(c *gin.Context, body *types.CreateCategoryParams) (*db.ChannelCategory, *types.APIError) {
	if ok := s.permissions.CheckPermission(c, body.ServerID, types.ManageChannels); !ok {
		return nil, types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN_CREATE_CATEGORY", "Forbidden to create category.", nil)
	}

	category, err := s.db.CreateCategory(c, body)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_CREATE_CATEGORY", "Failed to create category.", err)
	}

	return &category, nil
}

func (s *channelService) CreateChannel(c *gin.Context, body *types.CreateChannelParams) (*db.Channel, *types.APIError) {
	if ok := s.permissions.CheckPermission(c, body.ServerID, types.ManageChannels); !ok {
		return nil, types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN_CREATE_CHANNEL", "Forbidden to create channel.", nil)
	}

	channel, err := s.db.CreateChannel(c, body)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_CREATE_CHANNEL", "Failed to create channel.", err)
	}

	s.actors.StartChannel(body.ServerID, channel.ID)

	return &channel, nil
}

func (s *channelService) PinChannel(c *gin.Context, body *types.PinChannelParams) *types.APIError {
	user, exists := c.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Unauthorized.", nil)
	}
	userID := user.(*db.User).ID
	channelID := c.Param("channel_id")

	if err := s.db.PinChannel(c, channelID, userID, body); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_PIN_CHANNEL", "Failed to pin channel.", err)
	}

	return nil
}

func (s *channelService) DeleteChannel(c *gin.Context, body *types.DeleteChannelParams) *types.APIError {
	channelID := c.Param("channel_id")

	if ok := s.permissions.CheckPermission(c, body.ServerID, types.ManageChannels); !ok {
		return types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN_DELETE_CHANNEL", "Forbidden to delete channel.", nil)
	}

	if err := s.db.DeleteChannel(c, channelID); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_DELETE_CHANNEL", "Failed to delete channel.", err)
	}

	return nil
}

func (s *channelService) DeleteCategory(c *gin.Context, body *types.DeleteCategoryParams) *types.APIError {
	categoryID := c.Param("category_id")

	if ok := s.permissions.CheckPermission(c, body.ServerID, types.ManageChannels); !ok {
		return types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN_DELETE_CATEGORY", "Forbidden to delete category.", nil)
	}

	if err := s.db.DeleteCategory(c, categoryID); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_DELETE_CATEGORY", "Failed to delete category.", err)
	}

	return nil
}

func (s *channelService) EditChannel(c *gin.Context, body *types.EditChannelParams) *types.APIError {
	channelID := c.Param("channel_id")

	if ok := s.permissions.CheckPermission(c, body.ServerID, types.ManageChannels); !ok {
		return types.NewAPIError(http.StatusForbidden, "ERR_FORBIDDEN_DELETE_CATEGORY", "Forbidden to delete category.", nil)
	}

	if err := s.db.UpdateChannelInformations(c, channelID, body); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_EDIT_CHANNEL", "Failed to edit channel.", err)
	}

	return nil
}
