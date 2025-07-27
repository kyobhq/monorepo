package actors

import (
	"log"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/cluster"
	"github.com/lxzan/gws"
)

type Engine int

const (
	UserEngine Engine = iota
	ServerEngine
)

type Service interface {
	// CreateUser use the users actor engine to create a new user
	CreateUser(userID string, wsConn *gws.Conn) *actor.PID

	// GetUser use the users actor engine to get a user by id
	GetUser(userID string) *actor.PID

	// GetServer use the servers actor engine to get a server by id
	GetServer(serverID string) *actor.PID

	// GetChannel use the servers actor engine to get a channel by id
	GetChannel(serverID, channelID string) *actor.PID

	// SendMessageTo use any actor engine to send a message
	SendMessageTo(engine Engine, PID *actor.PID, mess any, userPID ...*actor.PID)

	KillActor(userPID *actor.PID)
}

type service struct {
	cluster *cluster.Cluster
}

func GetIDFromPID(PID *actor.PID) string {
	split := strings.Split(PID.ID, "/")
	return split[len(split)-1]
}

func New() Service {
	config := cluster.NewConfig().WithID("A").WithRegion("eu-west")
	c, err := cluster.New(config)
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}

	actorService := &service{
		cluster: c,
	}

	c.RegisterKind("server", newServer(actorService), cluster.NewKindConfig())
	c.RegisterKind("user", newUser(actorService, nil), cluster.NewKindConfig())
	c.Start()

	c.Spawn(newServer(actorService), "server", actor.WithID("global"))

	return actorService
}

func (se *service) CreateUser(userID string, wsConn *gws.Conn) *actor.PID {
	return se.cluster.Spawn(newUser(se, wsConn), "user", actor.WithID(userID))
}

func (se *service) GetUser(userID string) *actor.PID {
	return se.cluster.GetActiveByID("user/" + userID)
}

func (se *service) GetServer(serverID string) *actor.PID {
	return se.cluster.GetActiveByID("server/" + serverID)
}

func (se *service) GetChannel(serverID, channelID string) *actor.PID {
	serverPID := se.GetServer(serverID)
	return serverPID.Child("channel/" + channelID)
}

func (se *service) SendMessageTo(engine Engine, PID *actor.PID, mess any, userPID ...*actor.PID) {
	switch engine {
	case UserEngine:
		se.cluster.Engine().Send(PID, mess)
	case ServerEngine:
		se.cluster.Engine().SendWithSender(PID, mess, userPID[0])
	}
}

func (se *service) KillActor(userPID *actor.PID) {
	se.cluster.Deactivate(userPID)
}
