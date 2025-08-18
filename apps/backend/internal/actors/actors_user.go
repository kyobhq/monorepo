package actors

import (
	db "backend/db/gen_queries"
	"backend/internal/database"
	messages "backend/proto"
	"fmt"
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
	case *messages.GetFriends:
		fmt.Println("called")
		ctx.Respond(&messages.GetFriends{
			FriendIds: u.friends,
		})
	case *messages.AccountDeletion:
		u.AccountDeletion(ctx, msg)
	case *messages.ChangeStatus:
		u.FriendChangeStatus(ctx, msg)
	case *messages.WSMessage:
		message, _ := proto.Marshal(msg)
		u.wsConn.WriteMessage(gws.OpcodeBinary, message)
	}
}

func (u *user) FriendChangeStatus(ctx *actor.Context, msg *messages.ChangeStatus) {
	u.friends = append(u.friends, msg.User.Id)

	if msg.Type == "Ping" {
		userPID := u.hub.GetUser(msg.User.Id)
		ctx.Send(userPID, &messages.WSMessage{
			Content: &messages.WSMessage_UserChangeStatus{
				UserChangeStatus: &messages.ChangeStatus{
					Type: "connect",
					User: &messages.User{
						Id: GetIDFromPID(ctx.PID()),
					},
					Status: "online",
				},
			},
		})
	}

	m := &messages.WSMessage{
		Content: &messages.WSMessage_UserChangeStatus{
			UserChangeStatus: &messages.ChangeStatus{
				Type:   "connect",
				User:   msg.User,
				Status: msg.Status,
			},
		},
	}
	message, _ := proto.Marshal(m)
	u.wsConn.WriteMessage(gws.OpcodeBinary, message)
}

func (u *user) AccountDeletion(ctx *actor.Context, msg *messages.AccountDeletion) {
	userID := GetIDFromPID(ctx.PID())
	if msg.UserId == userID {
		for _, friendID := range u.friends {
			userPID := u.hub.GetUser(friendID)
			ctx.Send(userPID, msg)
		}
	} else {
		u.friends = slices.DeleteFunc(u.friends, func(friendID string) bool {
			return friendID == msg.UserId
		})

		m := &messages.WSMessage{
			Content: &messages.WSMessage_AccountDeletion{
				AccountDeletion: msg,
			},
		}

		message, _ := proto.Marshal(m)
		u.wsConn.WriteMessage(gws.OpcodeBinary, message)
	}
}

func (u *user) initializeUser(ctx *actor.Context) {
	userID := GetIDFromPID(ctx.PID())
	user, err := u.db.GetUserByID(ctx.Context(), userID)
	if err != nil {
		slog.Error("failed to get user", "err", err)
	}

	friendIDs, err := u.db.GetFriendIDs(ctx.Context(), userID)
	if err != nil {
		slog.Error("failed to get friendIDs", "err", err)
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

	for _, friendID := range friendIDs {
		connectMessage := &messages.ChangeStatus{
			Type: "Ping",
			User: &messages.User{
				Id: user.ID,
			},
			Status: "online",
		}

		u.hub.NotifyFriendStatus(friendID, connectMessage)
	}
}

func (u *user) killUser(ctx *actor.Context) {
	userID := GetIDFromPID(ctx.PID())
	servers, err := u.db.GetUserServers(ctx.Context(), userID)
	if err != nil {
		slog.Error("failed to get servers on disconnect", "err", err)
	}

	friendIDs, err := u.db.GetFriendIDs(ctx.Context(), userID)
	if err != nil {
		slog.Error("failed to get friends on disconnect", "err", err)
	}

	for _, server := range servers {
		disconnectMessage := &messages.ChangeStatus{
			Type: "disconnect",
			User: &messages.User{
				Id: userID,
			},
			ServerId: server.ID,
			Status:   "offline",
		}

		u.hub.SendUserStatusMessage(ctx.PID(), disconnectMessage)
	}

	for _, friendID := range friendIDs {
		disconnectMessage := &messages.ChangeStatus{
			Type: "disconnect",
			User: &messages.User{
				Id: userID,
			},
			Status: "offline",
		}

		u.hub.NotifyFriendStatus(friendID, disconnectMessage)
	}
}
