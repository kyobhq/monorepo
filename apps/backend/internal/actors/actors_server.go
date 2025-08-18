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
	case *messages.AccountDeletion:
		s.AccountDeletion(ctx, msg)
	case *messages.StartCategory:
		s.startCategory(msg)
	case *messages.KillCategory:
		s.killCategory(ctx, msg)
	case *messages.StartChannel:
		s.startChannel(ctx, msg)
	case *messages.KillChannel:
		s.killChannel(ctx, msg)
	case *messages.CreateOrEditRole:
		s.CreateOrEditRole(msg)
	case *messages.RemoveRole:
		s.RemoveRole(msg)
	case *messages.MoveRole:
		s.MoveRole(msg)
	case *messages.AddRoleMember:
		s.AddRoleMember(msg)
	case *messages.RemoveRoleMember:
		s.RemoveRoleMember(msg)
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

func (s *server) startCategory(msg *messages.StartCategory) {
	message := &messages.WSMessage{
		Content: &messages.WSMessage_StartCategory{
			StartCategory: msg,
		},
	}

	for userID := range s.users {
		s.hub.BroadcastMessageToUser(s.hub.GetUser(userID), message)
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

func (s *server) CreateOrEditRole(msg *messages.CreateOrEditRole) {
	message := &messages.WSMessage{
		Content: &messages.WSMessage_CreateOrEditRole{
			CreateOrEditRole: msg,
		},
	}

	for userID := range s.users {
		userPID := s.hub.GetUser(userID)
		s.hub.BroadcastMessageToUser(userPID, message)
	}
}

func (s *server) RemoveRole(msg *messages.RemoveRole) {
	message := &messages.WSMessage{
		Content: &messages.WSMessage_RemoveRole{
			RemoveRole: msg,
		},
	}

	for userID := range s.users {
		userPID := s.hub.GetUser(userID)
		s.hub.BroadcastMessageToUser(userPID, message)
	}
}

func (s *server) MoveRole(msg *messages.MoveRole) {
	message := &messages.WSMessage{
		Content: &messages.WSMessage_MoveRole{
			MoveRole: msg,
		},
	}

	for userID := range s.users {
		userPID := s.hub.GetUser(userID)
		s.hub.BroadcastMessageToUser(userPID, message)
	}
}

func (s *server) AddRoleMember(msg *messages.AddRoleMember) {
	message := &messages.WSMessage{
		Content: &messages.WSMessage_AddRoleMember{
			AddRoleMember: msg,
		},
	}

	for userID := range s.users {
		userPID := s.hub.GetUser(userID)
		s.hub.BroadcastMessageToUser(userPID, message)
	}
}

func (s *server) RemoveRoleMember(msg *messages.RemoveRoleMember) {
	message := &messages.WSMessage{
		Content: &messages.WSMessage_RemoveRoleMember{
			RemoveRoleMember: msg,
		},
	}

	for userID := range s.users {
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
