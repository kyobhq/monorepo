package actors

import (
	"backend/internal/database"
	protoTypes "backend/proto"
	"context"
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
	CreateUser(userID string, wsConn *gws.Conn) *actor.PID

	GetUser(userID string) *actor.PID

	// GetAllServersInstances(serverIDs []string) []*actor.PID

	GetAllServerInstances(serverID string) []*actor.PID

	GetAllChannelInstances(serverID, channelID string) []*actor.PID

	SendMessageTo(engine Engine, PID *actor.PID, mess any, userPID ...*actor.PID)

	KillActor(userPID *actor.PID)

	StartServerInRegion(serverID, region string) *actor.PID
}

type service struct {
	cluster *cluster.Cluster
	db      database.Service
}

func GetIDFromPID(PID *actor.PID) string {
	split := strings.Split(PID.ID, "/")
	return split[len(split)-1]
}

func New(dbService database.Service) Service {
	config := cluster.NewConfig().WithID(os.Getenv("NODE_ID")).WithRegion(os.Getenv("REGION")).WithListenAddr(os.Getenv("NODE_IP"))
	c, err := cluster.New(config)
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}

	actorService := &service{
		cluster: c,
		db:      dbService,
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
	actorService.Bootstrap()

	return actorService
}

func (se *service) Bootstrap() {
	serverIDs, err1 := se.db.GetServers(context.TODO())
	channels, err2 := se.db.GetChannels(context.TODO())
	if err1 != nil || err2 != nil {
		log.Fatal(err1, err2)
	}

	nodeRegion := os.Getenv("REGION")
	for _, serverID := range serverIDs {
		if actorPID := se.cluster.GetActiveByID("server/" + serverID + "@" + nodeRegion); actorPID == nil {
			serverPID := se.StartServerInRegion(serverID, nodeRegion)

			for _, channel := range channels {
				if channel.ServerID == serverID {
					se.cluster.Engine().Send(serverPID, protoTypes.StartChannel{
						Channel: &protoTypes.Channel{
							Id: channel.ID,
						},
					})
				}
			}
		}
	}
}

func (se *service) CreateUser(userID string, wsConn *gws.Conn) *actor.PID {
	return se.cluster.Spawn(newUser(se, wsConn), "user", actor.WithID(userID))
}

func (se *service) GetUser(userID string) *actor.PID {
	return se.cluster.GetActiveByID("user/" + userID)
}

// func (se *service) GetAllServersInstances(serverIDs []string) []*actor.PID {
// 	var instances []*actor.PID
//
// 	for _, id := range serverIDs {
// 		for _, region := range regions {
// 			actorPID := se.cluster.GetActiveByID("server/" + id + "@" + region)
// 			if actorPID == nil {
// 				actorPID = se.StartServerInRegion(id, region)
// 			}
//
// 			if region == os.Getenv("region") {
// 				instances = append(instances, actorPID)
// 			}
// 		}
// 	}
//
// 	return instances
// }

func (se *service) GetAllServerInstances(serverID string) []*actor.PID {
	var instances []*actor.PID

	for _, region := range regions {
		actorPID := se.cluster.GetActiveByID("server/" + serverID + "@" + region)
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
	return se.cluster.Activate("server", cluster.NewActivationConfig().WithID(serverID+"@"+region).WithRegion(region))
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
