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
		if len(c.users) == 0 {
			response := ctx.Request(ctx.Parent(), &messages.GetServerUsers{}, 10*time.Second)
			result, err := response.Result()
			if err == nil {
				c.NewMessage(ctx, result.(*messages.GetServerUsers).UserIds, msg)
			}
		} else {
			c.NewMessage(ctx, c.users, msg)
		}
	}
}

func (c *channel) NewMessage(ctx *actor.Context, userIDs []string, msg *messages.NewChatMessage) {
	messageToBroadcast := &messages.WSMessage{
		Type: "chat_message",
		Content: &messages.WSMessage_NewChatMessage{
			NewChatMessage: msg,
		},
	}

	for _, userID := range userIDs {
		userPID := c.hub.GetUser(userID)
		c.hub.SendMessageTo(UserEngine, userPID, messageToBroadcast)
	}
}
