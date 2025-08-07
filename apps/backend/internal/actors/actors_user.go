package actors

import (
	db "backend/db/gen_queries"
	"backend/internal/database"
	messages "backend/proto"
	"log/slog"
	"slices"

	"github.com/anthdm/hollywood/actor"
	"github.com/lxzan/gws"
	"google.golang.org/protobuf/proto"
)

type user struct {
	logger  *slog.Logger
	wsConn  *gws.Conn
	friends []string
	hub     Service
	db      database.Service
}

func newUser(actorService Service, db database.Service, wsConn *gws.Conn) actor.Producer {
	return func() actor.Receiver {
		return &user{
			logger:  slog.Default(),
			wsConn:  wsConn,
			friends: []string{},
			hub:     actorService,
			db:      db,
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
	}
}

func (u *user) initializeUser(ctx *actor.Context) {
	userID := GetIDFromPID(ctx.PID())
	user, err := u.db.GetUserByID(ctx.Context(), userID)
	if err != nil {
		slog.Error("failed to get user", "err", err)
	}

	serverIDs, err := u.db.GetServersIDFromUser(ctx.Context(), userID)
	if err != nil {
		slog.Error("failed to get serverIDs", "err", err)
	}

	roles, err := u.db.GetUserRolesFromServers(ctx.Context(), userID, serverIDs)
	if err != nil {
		slog.Error("failed to get roles", "err", err)
	}

	for _, serverID := range serverIDs {
		idx := slices.IndexFunc(roles, func(s db.GetUserRolesFromServersRow) bool {
			return s.ServerID == serverID
		})

		connectMessage := &messages.ChangeStatus{
			Type: "connect",
			User: &messages.User{
				Id:          user.ID,
				DisplayName: user.DisplayName,
				Avatar:      user.Avatar.String,
			},
			ServerId: serverID,
			Status:   "online",
			Roles:    roles[idx].Roles,
		}

		u.hub.SendUserStatusMessage(ctx.PID(), connectMessage)
	}
}

func (u *user) killUser(ctx *actor.Context) {
	userID := GetIDFromPID(ctx.PID())
	servers, err := u.db.GetUserServers(ctx.Context(), userID)
	if err != nil {
		slog.Error("failed to initialize user", "err", err)
	}

	for _, server := range servers {
		disconnectMessage := &messages.ChangeStatus{
			User: &messages.User{
				Id: userID,
			},
			ServerId: server.ID,
			Status:   "offline",
		}

		u.hub.SendUserStatusMessage(ctx.PID(), disconnectMessage)
	}
}
