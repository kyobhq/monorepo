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
		slog.Info("server started",
			"pid", ctx.PID(),
		)
		ctx.SpawnChild(newChannel(s.hub), "channel", actor.WithID(msg.Channel.Id))
	case *messages.GetServerUsers:
		ctx.Respond(&messages.GetServerUsers{
			UserIds: slices.Collect(maps.Keys(s.users)),
		})
	case *messages.ChangeStatus:
		switch msg.Status {
		case "online":
			s.users[msg.Id] = Online
		case "offline":
			delete(s.users, msg.Id)
		}
	}
}
