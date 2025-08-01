package actors

import (
	"fmt"
	"log"
	"os"
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

var regions = []string{"na", "eu", "asia"}

type Service interface {
	// CreateUser use the users actor engine to create a new user
	CreateUser(userID string, wsConn *gws.Conn) *actor.PID

	// GetUser use the users actor engine to get a user by id
	GetUser(userID string) *actor.PID

	// GetServer use the servers actor engine to get a server by id
	GetAllServerInstances(serverID string) []*actor.PID

	// GetChannel use the servers actor engine to get a channel by id
	GetAllChannelInstances(serverID, channelID string) []*actor.PID

	// SendMessageTo use any actor engine to send a message
	SendMessageTo(engine Engine, PID *actor.PID, mess any, userPID ...*actor.PID)

	KillActor(userPID *actor.PID)

	StartServerInRegion(serverID, region string) *actor.PID
}

type service struct {
	cluster *cluster.Cluster
}

func GetIDFromPID(PID *actor.PID) string {
	split := strings.Split(PID.ID, "/")
	return split[len(split)-1]
}

func New() Service {
	config := cluster.NewConfig().WithID(os.Getenv("NODE_ID")).WithRegion(os.Getenv("REGION")).WithListenAddr(os.Getenv("NODE_IP"))
	c, err := cluster.New(config)
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}

	actorService := &service{
		cluster: c,
	}

	c.RegisterKind("server", newServer(actorService), cluster.NewKindConfig())
	c.RegisterKind("user", newUser(actorService, nil), cluster.NewKindConfig())

	eventPID := c.Engine().SpawnFunc(func(ctx *actor.Context) {
		switch msg := ctx.Message().(type) {
		case cluster.ActivationEvent:
			fmt.Println("got activation event")
		case cluster.MemberJoinEvent:
			fmt.Println("member joined", msg.Member.Host, msg.Member.ID, msg.Member.Region)
		}
	}, "event")
	c.Engine().Subscribe(eventPID)

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

func (se *service) GetAllServerInstances(serverID string) []*actor.PID {
	var instances []*actor.PID

	for _, region := range regions {
		actorPID := se.cluster.GetActiveByID("server/" + serverID + "@" + region)
		if actorPID == nil {
			actorPID = se.StartServerInRegion(serverID, region)
		}

		instances = append(instances, actorPID)
	}

	return instances
}

func (se *service) GetAllChannelInstances(serverID, channelID string) []*actor.PID {
	var instances []*actor.PID

	serverPIDs := se.GetAllServerInstances(serverID)
	for _, pid := range serverPIDs {
		instances = append(instances, pid.Child("channel/"+channelID))
	}

	return instances
}

func (se *service) StartServerInRegion(serverID, region string) *actor.PID {
	return se.cluster.Activate("server", cluster.NewActivationConfig().WithID(serverID).WithRegion(region))
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
