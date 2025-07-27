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
	serverPID := u.hub.GetServer("global")
	u.hub.SendMessageTo(ServerEngine, serverPID, connectMessage, ctx.PID())
}

func (u *user) killUser(ctx *actor.Context) {
	disconnectMessage := &messages.ChangeStatus{
		Id:       GetIDFromPID(ctx.PID()),
		ServerId: "global",
		Status:   "offline",
	}
	serverPID := u.hub.GetServer("global")
	u.hub.SendMessageTo(ServerEngine, serverPID, disconnectMessage, ctx.PID())
}
