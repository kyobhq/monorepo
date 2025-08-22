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
	case actor.InternalError:
		slog.Error("server erroring",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	case *messages.WSMessage:
		s.BroadcastToServer(msg)
	case *messages.AccountDeletion:
		s.AccountDeletion(ctx, msg)
	case *messages.KillCategory:
		s.killCategory(ctx, msg)
	case *messages.StartChannel:
		s.startChannel(ctx, msg)
	case *messages.KillChannel:
		s.killChannel(ctx, msg)
	case *messages.LeaveServer:
		s.LeaveServer(msg)
	case *messages.KickUser:
		s.KickUser(msg)
	case *messages.BanUser:
		s.BanUser(msg)
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
	ctx.SpawnChild(newChannel(s.hub, msg.Channel.Users), "channel", actor.WithID(msg.Channel.Id))

	if msg.Channel.ServerId == "global" {
		return
	}

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

func (s *server) killCategory(ctx *actor.Context, msg *messages.KillCategory) {
	for _, channelID := range msg.ChannelsIds {
		channelPID := ctx.PID().Child("channel/" + channelID)
		ctx.Engine().Poison(channelPID)
	}

	message := &messages.WSMessage{
		Content: &messages.WSMessage_KillCategory{
			KillCategory: &messages.KillCategory{
				ServerId:   msg.ServerId,
				CategoryId: msg.CategoryId,
			},
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

func (s *server) AccountDeletion(ctx *actor.Context, msg *messages.AccountDeletion) {
	delete(s.users, msg.UserId)

	for _, channelPID := range ctx.Children() {
		ctx.Send(channelPID, msg)
	}
}

func (s *server) BroadcastToServer(msg *messages.WSMessage) {
	for userID := range s.users {
		userPID := s.hub.GetUser(userID)
		s.hub.BroadcastMessageToUser(userPID, msg)
	}
}

func (s *server) LeaveServer(msg *messages.LeaveServer) {
	delete(s.users, msg.UserId)

	message := &messages.WSMessage{
		Content: &messages.WSMessage_LeaveServer{
			LeaveServer: msg,
		},
	}

	for userID := range s.users {
		userPID := s.hub.GetUser(userID)
		s.hub.BroadcastMessageToUser(userPID, message)
	}
}

func (s *server) BanUser(msg *messages.BanUser) {
	message := &messages.WSMessage{
		Content: &messages.WSMessage_BanUser{
			BanUser: msg,
		},
	}

	for userID := range s.users {
		userPID := s.hub.GetUser(userID)
		s.hub.BroadcastMessageToUser(userPID, message)
	}

	delete(s.users, msg.UserId)
}

func (s *server) KickUser(msg *messages.KickUser) {
	message := &messages.WSMessage{
		Content: &messages.WSMessage_KickUser{
			KickUser: msg,
		},
	}

	for userID := range s.users {
		userPID := s.hub.GetUser(userID)
		s.hub.BroadcastMessageToUser(userPID, message)
	}

	delete(s.users, msg.UserId)
}
