package actors

import (
	messages "backend/proto"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type channel struct {
	logger *slog.Logger
	users  []string
	hub    Service
}

func newChannel(actorService Service) actor.Producer {
	return func() actor.Receiver {
		return &channel{
			logger: slog.Default(),
			users:  []string{},
			hub:    actorService,
		}
	}
}

func (c *channel) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slog.Info("channel started",
			"id", ctx.PID().GetID(),
		)
	case actor.Stopped:
		slog.Info("channel stopped",
			"id", ctx.PID().GetID(),
		)
	case actor.InternalError:
		slog.Error("channel erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	case *messages.NewChatMessage:
		c.NewMessage(ctx, c.GetChannelUsers(ctx), msg)
	case *messages.EditChatMessage:
		c.EditMessage(ctx, c.GetChannelUsers(ctx), msg)
	case *messages.DeleteChatMessage:
		c.DeleteMessage(ctx, c.GetChannelUsers(ctx), msg)
	}
}

func (c *channel) GetChannelUsers(ctx *actor.Context) []string {
	if len(c.users) == 0 {
		response := ctx.Request(ctx.Parent(), &messages.GetServerUsers{}, 10*time.Second)
		result, err := response.Result()
		if err == nil {
			return result.(*messages.GetServerUsers).UserIds
		}
	}

	return c.users
}

func (c *channel) NewMessage(ctx *actor.Context, userIDs []string, msg *messages.NewChatMessage) {
	messageToBroadcast := &messages.WSMessage{
		Content: &messages.WSMessage_NewChatMessage{
			NewChatMessage: msg,
		},
	}

	for _, userID := range userIDs {
		userPID := c.hub.GetUser(userID)
		c.hub.BroadcastMessageToUser(userPID, messageToBroadcast)
	}
}

func (c *channel) EditMessage(ctx *actor.Context, userIDs []string, msg *messages.EditChatMessage) {
	messageToBroadcast := &messages.WSMessage{
		Content: &messages.WSMessage_EditChatMessage{
			EditChatMessage: msg,
		},
	}

	for _, userID := range userIDs {
		userPID := c.hub.GetUser(userID)
		c.hub.BroadcastMessageToUser(userPID, messageToBroadcast)
	}
}

func (c *channel) DeleteMessage(ctx *actor.Context, userIDs []string, msg *messages.DeleteChatMessage) {
	messageToBroadcast := &messages.WSMessage{
		Content: &messages.WSMessage_DeleteChatMessage{
			DeleteChatMessage: msg,
		},
	}

	for _, userID := range userIDs {
		userPID := c.hub.GetUser(userID)
		c.hub.BroadcastMessageToUser(userPID, messageToBroadcast)
	}
}
