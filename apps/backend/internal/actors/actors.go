package actors

import (
	db "backend/db/gen_queries"
	"backend/internal/database"
	"backend/internal/types"
	message "backend/proto"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/cluster"
	"github.com/lxzan/gws"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	SendChatMessage(chatMessage *message.NewChatMessage)

	EditMessage(chatMessage *message.EditChatMessage)

	DeleteMessage(chatMessage *message.DeleteChatMessage)

	SendUserStatusMessage(userPID *actor.PID, status *message.ChangeStatus)

	KillActor(userPID *actor.PID)

	StartServerInRegion(serverID, region string) *actor.PID

	StartCategory(category db.ChannelCategory)

	StartChannel(channel db.Channel)

	KillChannel(body *types.DeleteChannelParams, channelID string)

	KillCategory(body *types.DeleteCategoryParams, categoryID string)

	BroadcastMessageToUser(userPID *actor.PID, message *message.WSMessage)

	GetActiveUsers(serverID string) []string
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
	c.RegisterKind("user", newUser(actorService, dbService, nil), cluster.NewKindConfig())

	eventPID := c.Engine().SpawnFunc(func(ctx *actor.Context) {
		switch msg := ctx.Message().(type) {
		case cluster.ActivationEvent:
			// fmt.Println("got activation event")
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
		serverPID := se.StartServerInRegion(serverID, nodeRegion)

		for _, channel := range channels {
			if channel.ServerID == serverID {
				se.cluster.Engine().Send(serverPID, &message.StartChannel{
					Channel: &message.Channel{
						Id: channel.ID,
					},
				})
			}
		}
	}
}

func (se *service) KillActor(userPID *actor.PID) {
	se.cluster.Deactivate(userPID)
}

func (se *service) CreateUser(userID string, wsConn *gws.Conn) *actor.PID {
	return se.cluster.Spawn(newUser(se, se.db, wsConn), "user", actor.WithID(userID))
}

func (se *service) GetUser(userID string) *actor.PID {
	return se.cluster.GetActiveByID("user/" + userID)
}

func (se *service) GetAllServerInstances(serverID string) []*actor.PID {
	var instances []*actor.PID

	for _, region := range regions {
		actorPID := se.cluster.GetActiveByID("server/" + serverID + "@" + region)

		if actorPID != nil {
			instances = append(instances, actorPID)
		}
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
	actorID := "server/" + serverID + "@" + region
	if existingPID := se.cluster.GetActiveByID(actorID); existingPID != nil {
		return existingPID
	}

	return se.cluster.Activate("server", cluster.NewActivationConfig().WithID(serverID+"@"+region).WithRegion(region))
}

func (se *service) StartCategory(category db.ChannelCategory) {
	serversPID := se.GetAllServerInstances(category.ServerID)

	for _, serverPID := range serversPID {
		se.cluster.Engine().Send(serverPID, &message.StartCategory{
			Category: &message.Category{
				Id:        category.ID,
				Position:  category.Position,
				ServerId:  category.ServerID,
				Name:      category.Name,
				Users:     category.Users,
				Roles:     category.Roles,
				E2Ee:      category.E2ee,
				CreatedAt: timestamppb.New(category.CreatedAt),
				UpdatedAt: timestamppb.New(category.UpdatedAt),
			},
		})
	}
}

func (se *service) StartChannel(channel db.Channel) {
	serversPID := se.GetAllServerInstances(channel.ServerID)

	for _, serverPID := range serversPID {
		se.cluster.Engine().Send(serverPID, &message.StartChannel{
			Channel: &message.Channel{
				Id:          channel.ID,
				ServerId:    channel.ServerID,
				CategoryId:  channel.CategoryID,
				Name:        channel.Name,
				Description: channel.Description.String,
				Type:        channel.Type,
				E2Ee:        channel.E2ee,
				Users:       channel.Users,
				Roles:       channel.Roles,
				Position:    channel.Position,
				Active:      channel.Active,
				CreatedAt:   timestamppb.New(channel.CreatedAt),
				UpdatedAt:   timestamppb.New(channel.UpdatedAt),
			},
		})
	}
}

func (se *service) KillCategory(body *types.DeleteCategoryParams, categoryID string) {
	serversPID := se.GetAllServerInstances(body.ServerID)

	for _, serverPID := range serversPID {
		se.cluster.Engine().Send(serverPID, &message.KillCategory{
			ServerId:    body.ServerID,
			CategoryId:  categoryID,
			ChannelsIds: body.ChannelsIDs,
		})
	}
}

func (se *service) KillChannel(body *types.DeleteChannelParams, channelID string) {
	serversPID := se.GetAllServerInstances(body.ServerID)

	for _, serverPID := range serversPID {
		se.cluster.Engine().Send(serverPID, &message.KillChannel{
			Channel: &message.Channel{
				Id:         channelID,
				ServerId:   body.ServerID,
				CategoryId: body.CategoryID,
			},
		})
	}
}

func (se *service) SendChatMessage(chatMessage *message.NewChatMessage) {
	channels := se.GetAllChannelInstances(chatMessage.Message.ServerId, chatMessage.Message.ChannelId)
	for _, channelPID := range channels {
		se.cluster.Engine().Send(channelPID, chatMessage)
	}
}

func (se *service) EditMessage(chatMessage *message.EditChatMessage) {
	channels := se.GetAllChannelInstances(chatMessage.Message.ServerId, chatMessage.Message.ChannelId)
	for _, channelPID := range channels {
		se.cluster.Engine().Send(channelPID, chatMessage)
	}
}

func (se *service) DeleteMessage(chatMessage *message.DeleteChatMessage) {
	channels := se.GetAllChannelInstances(chatMessage.Message.ServerId, chatMessage.Message.ChannelId)
	for _, channelPID := range channels {
		se.cluster.Engine().Send(channelPID, chatMessage)
	}
}

func (se *service) SendUserStatusMessage(userPID *actor.PID, status *message.ChangeStatus) {
	servers := se.GetAllServerInstances(status.ServerId)
	for _, serverPID := range servers {
		se.cluster.Engine().SendWithSender(serverPID, status, userPID)
	}
}

func (se *service) BroadcastMessageToUser(userPID *actor.PID, message *message.WSMessage) {
	se.cluster.Engine().Send(userPID, message)
}

func (se *service) GetActiveUsers(serverID string) []string {
	var allUsersIDs []string

	servers := se.GetAllServerInstances(serverID)
	for _, server := range servers {
		response := se.cluster.Engine().Request(server, &message.GetServerUsers{}, 10*time.Second)
		result, err := response.Result()
		if err == nil {
			allUsersIDs = append(allUsersIDs, result.(*message.GetServerUsers).UserIds...)
		}
	}

	return allUsersIDs
}
