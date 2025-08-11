package actors

import (
	messages "backend/proto"
	"log/slog"
	"maps"
	"slices"

	"github.com/anthdm/hollywood/actor"
)

type Status string

const (
	Online Status = "online"
	Dnd    Status = "dnd"
	Away   Status = "away"
)

type server struct {
	logger *slog.Logger
	users  map[string]Status
	hub    Service
}

func newServer(actorService Service) actor.Producer {
	return func() actor.Receiver {
		return &server{
			logger: slog.Default(),
			users:  make(map[string]Status),
			hub:    actorService,
		}
	}
}

func (s *server) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slog.Info("server started",
			"pid", ctx.PID(),
		)
	case actor.Stopped:
		slog.Info("server stopped",
			"id", ctx.PID().GetID(),
		)
	case actor.InternalError:
		slog.Error("server erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	case *messages.StartChannel:
		s.startChannel(ctx, msg)
	case *messages.KillChannel:
		s.killChannel(ctx, msg)
	case *messages.GetServerUsers:
		ctx.Respond(&messages.GetServerUsers{
			UserIds: slices.Collect(maps.Keys(s.users)),
		})
	case *messages.ChangeStatus:
		s.broadcastUserStatus(msg)
		switch msg.Status {
		case "online":
			s.users[msg.User.Id] = Online
		case "offline":
			delete(s.users, msg.User.Id)
		}
	}
}

func (s *server) startChannel(ctx *actor.Context, msg *messages.StartChannel) {
	ctx.SpawnChild(newChannel(s.hub), "channel", actor.WithID(msg.Channel.Id))

	message := &messages.WSMessage{
		Content: &messages.WSMessage_StartChannel{
			StartChannel: msg,
		},
	}

	for userID := range s.users {
		s.hub.BroadcastMessageToUser(s.hub.GetUser(userID), message)
	}
}

func (s *server) killChannel(ctx *actor.Context, msg *messages.KillChannel) {
	channelPID := ctx.PID().Child("channel/" + msg.Channel.Id)
	ctx.Engine().Poison(channelPID)

	message := &messages.WSMessage{
		Content: &messages.WSMessage_KillChannel{
			KillChannel: msg,
		},
	}

	for userID := range s.users {
		s.hub.BroadcastMessageToUser(s.hub.GetUser(userID), message)
	}
}

func (s *server) broadcastUserStatus(msg *messages.ChangeStatus) {
	message := &messages.WSMessage{
		Content: &messages.WSMessage_UserChangeStatus{
			UserChangeStatus: msg,
		},
	}

	for userID := range s.users {
		if userID == msg.User.Id {
			continue
		}

		userPID := s.hub.GetUser(userID)
		s.hub.BroadcastMessageToUser(userPID, message)
	}
}
