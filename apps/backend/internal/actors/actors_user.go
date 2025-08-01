package actors

import (
	messages "backend/proto"
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	"github.com/lxzan/gws"
	"google.golang.org/protobuf/proto"
)

type user struct {
	logger  *slog.Logger
	wsConn  *gws.Conn
	friends []string
	hub     Service
}

func newUser(actorService Service, wsConn *gws.Conn) actor.Producer {
	return func() actor.Receiver {
		return &user{
			logger:  slog.Default(),
			wsConn:  wsConn,
			friends: []string{},
			hub:     actorService,
		}
	}
}

func (u *user) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		u.initializeUser(ctx)
	case actor.Stopped:
		u.killUser(ctx)
	case actor.InternalError:
		slog.Error("actor user internal error",
			"id", ctx.PID().GetID(),
			"err", msg.Err,
		)
	case *messages.WSMessage:
		message, _ := proto.Marshal(msg)
		u.wsConn.WriteMessage(gws.OpcodeBinary, message)
	default:
	}
}

func (u *user) initializeUser(ctx *actor.Context) {
	connectMessage := &messages.ChangeStatus{
		Id:        GetIDFromPID(ctx.PID()),
		ServerId:  "global",
		ChannelId: "global",
		Status:    "online",
	}

	serverPIDs := u.hub.GetAllServerInstances("global")
	for _, pid := range serverPIDs {
		u.hub.SendMessageTo(ServerEngine, pid, connectMessage, ctx.PID())
	}
}

func (u *user) killUser(ctx *actor.Context) {
	disconnectMessage := &messages.ChangeStatus{
		Id:       GetIDFromPID(ctx.PID()),
		ServerId: "global",
		Status:   "offline",
	}

	serverPIDs := u.hub.GetAllServerInstances("global")
	for _, pid := range serverPIDs {
		u.hub.SendMessageTo(ServerEngine, pid, disconnectMessage, ctx.PID())
	}
}
