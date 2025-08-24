package domains

import (
	db "backend/db/gen_queries"
	"backend/internal/actors"
	"backend/internal/database"
	"backend/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FriendService interface {
	SendRequest(ctx *gin.Context, body *types.SendRequestParams) (*db.GetFriendsRow, *types.APIError)
	AcceptRequest(ctx *gin.Context, body *types.AcceptRequestParams) *types.APIError
	RemoveFriend(ctx *gin.Context, body *types.RemoveFriendParams) *types.APIError
}

type friendService struct {
	db     database.Service
	actors actors.Service
}

func NewFriendService(db database.Service, actors actors.Service) *friendService {
	return &friendService{
		db:     db,
		actors: actors,
	}
}

func (s *friendService) SendRequest(ctx *gin.Context, body *types.SendRequestParams) (*db.GetFriendsRow, *types.APIError) {
	u, exists := ctx.Get("user")
	if !exists {
		return nil, types.NewAPIError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
	}
	user := u.(*db.User)

	receiver, err := s.db.GetUser(ctx, body.ReceiverUsername)
	if err != nil {
		return nil, types.NewAPIError(http.StatusNotFound, "ERR_USER_NOT_FOUND", "User not found", err)
	}

	friendship, err := s.db.CreateFriendRequest(ctx, user.ID, receiver.ID)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_FRIENDSHIP_CREATION", "Failed to create friend request", err)
	}

	s.actors.SendFriendRequest(friendship.ID, receiver.ID, user)

	return &db.GetFriendsRow{
		ID:                 receiver.ID,
		DisplayName:        receiver.DisplayName,
		Avatar:             receiver.Avatar,
		Banner:             receiver.Banner,
		AboutMe:            receiver.AboutMe,
		Accepted:           friendship.Accepted,
		FriendshipID:       friendship.ID,
		FriendshipSenderID: user.ID,
	}, nil
}

func (s *friendService) AcceptRequest(ctx *gin.Context, body *types.AcceptRequestParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
	}
	user := u.(*db.User)

	channelID, err := s.db.AcceptFriendRequest(ctx, body.FriendshipID, body.SenderID, user.ID)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_FRIENDSHIP_ACCEPT", "Failed to accept friend request", err)
	}

	s.actors.AcceptFriendRequest(body.FriendshipID, body.SenderID, user.ID, *channelID)

	return nil
}

func (s *friendService) RemoveFriend(ctx *gin.Context, body *types.RemoveFriendParams) *types.APIError {
	u, exists := ctx.Get("user")
	if !exists {
		return types.NewAPIError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", nil)
	}
	user := u.(*db.User)

	if err := s.db.RemoveFriend(ctx, body.FriendshipID, user.ID); err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_FRIENDSHIP_REMOVE", "Failed to remove friend", err)
	}

	if body.ChannelID != "" {
		s.actors.RemoveFriend(body.FriendshipID, body.SenderID, body.ReceiverID, body.ChannelID)
	}

	return nil
}
