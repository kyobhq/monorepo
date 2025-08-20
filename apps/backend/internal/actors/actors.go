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

	StartDMChannel(channelID string, userIDs []string)

	KillChannel(body *types.DeleteChannelParams, channelID string)

	KillCategory(body *types.DeleteCategoryParams, categoryID string)

	KillServer(serverID string)

	LeaveServer(serverID, userID string)

	BroadcastMessageToUser(userPID *actor.PID, message *message.WSMessage)

	GetActiveUsers(serverID string) []string

	CreateOrEditRole(role db.Role)

	RemoveRole(body *types.DeleteRoleParams)

	MoveRole(body *types.MoveRoleMemberParams)

	AddRoleMember(body *types.ChangeRoleMemberParams)

	RemoveRoleMember(body *types.ChangeRoleMemberParams)

	SendFriendRequest(friendshipID string, receiverID string, sender *db.User)

	AcceptFriendRequest(friendshipID, senderID, receiverID, channelID string)

	RemoveFriend(friendshipID, senderID, receiverID, channelID string)

	BanUser(serverID string, body *types.BanUserParams)

	KickUser(serverID string, body *types.KickUserParams)

	NotifyAccountDeletion(userID string, serverIDs []string)

	NotifyFriendStatus(friendID string, msg *message.ChangeStatus)

	GetActiveFriends(userID string) []string

	AvatarServerChange(serverID string, bannerURL, avatarURL *string)

	ProfileServerChange(serverID string, body *types.UpdateServerProfileParams)

	EditChannel(channelID string, body *types.EditChannelParams)

	EditCategory(categoryID string, body *types.EditCategoryParams)
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

	actorService.Bootstrap()

	return actorService
}

func (se *service) Bootstrap() {
	serverIDs, err1 := se.db.GetServers(context.TODO())
	channels, err2 := se.db.GetChannels(context.TODO())
	if err1 != nil || err2 != nil {
		log.Fatal(err1, err2)
	}

	// start servers
	nodeRegion := os.Getenv("REGION")
	for _, serverID := range serverIDs {
		serverPID := se.StartServerInRegion(serverID, nodeRegion)

		for _, channel := range channels {
			if channel.ServerID == serverID {
				se.cluster.Engine().Send(serverPID, &message.StartChannel{
					Channel: &message.Channel{
						Id:    channel.ID,
						Users: channel.Users,
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

func (se *service) NotifyFriendStatus(friendID string, msg *message.ChangeStatus) {
	friendPID := se.GetUser(friendID)
	se.cluster.Engine().Send(friendPID, msg)
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
		message := &message.WSMessage{
			Content: &message.WSMessage_StartCategory{
				StartCategory: &message.StartCategory{
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
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) StartChannel(channel db.Channel) {
	serversPID := se.GetAllServerInstances(channel.ServerID)

	for _, serverPID := range serversPID {
		se.cluster.Engine().Send(serverPID, &message.StartChannel{
			Channel: &message.Channel{
				Id:          channel.ID,
				ServerId:    channel.ServerID,
				CategoryId:  channel.CategoryID.String,
				Name:        channel.Name,
				Description: channel.Description.String,
				Type:        channel.Type,
				E2Ee:        channel.E2ee,
				Users:       channel.Users,
				Roles:       channel.Roles,
				Position:    channel.Position,
				CreatedAt:   timestamppb.New(channel.CreatedAt),
				UpdatedAt:   timestamppb.New(channel.UpdatedAt),
			},
		})
	}
}

func (se *service) StartDMChannel(channelID string, userIDs []string) {
	serversPID := se.GetAllServerInstances("global")

	for _, serverPID := range serversPID {
		se.cluster.Engine().Send(serverPID, &message.StartChannel{
			Channel: &message.Channel{
				Id:    channelID,
				Users: userIDs,
			},
		})
	}
}

func (se *service) KillServer(serverID string) {
	allUsers := se.GetActiveUsers(serverID)
	serversPID := se.GetAllServerInstances(serverID)

	for _, userID := range allUsers {
		userPID := se.GetUser(userID)

		message := &message.WSMessage{
			Content: &message.WSMessage_KillServer{
				KillServer: &message.KillServer{
					ServerId: serverID,
				},
			},
		}

		se.BroadcastMessageToUser(userPID, message)
	}

	for _, serverPID := range serversPID {
		se.cluster.Engine().Poison(serverPID)
	}
}

func (se *service) LeaveServer(serverID, userID string) {
	serversPID := se.GetAllServerInstances(serverID)

	for _, serverPID := range serversPID {
		message := &message.LeaveServer{
			ServerId: serverID,
			UserId:   userID,
		}

		se.cluster.Engine().Send(serverPID, message)
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

func (se *service) CreateOrEditRole(role db.Role) {
	serversPID := se.GetAllServerInstances(role.ServerID)

	for _, serverPID := range serversPID {
		message := &message.WSMessage{
			Content: &message.WSMessage_CreateOrEditRole{
				CreateOrEditRole: &message.CreateOrEditRole{
					Role: &message.Role{
						Id:        role.ID,
						ServerId:  role.ServerID,
						Position:  role.Position,
						Name:      role.Name,
						Color:     role.Color,
						Abilities: role.Abilities,
						CreatedAt: timestamppb.New(role.CreatedAt),
						UpdatedAt: timestamppb.New(role.UpdatedAt),
					},
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) RemoveRole(body *types.DeleteRoleParams) {
	serversPID := se.GetAllServerInstances(body.ServerID)

	for _, serverPID := range serversPID {
		message := &message.WSMessage{
			Content: &message.WSMessage_RemoveRole{
				RemoveRole: &message.RemoveRole{
					Role: &message.Role{
						Id:       body.RoleID,
						ServerId: body.ServerID,
					},
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) MoveRole(body *types.MoveRoleMemberParams) {
	serversPID := se.GetAllServerInstances(body.ServerID)

	for _, serverPID := range serversPID {
		message := &message.WSMessage{
			Content: &message.WSMessage_MoveRole{
				MoveRole: &message.MoveRole{
					MovedRole: &message.Role{
						Id:       body.MovedRoleID,
						ServerId: body.ServerID,
					},
					TargetRole: &message.Role{
						Id:       body.TargetRoleID,
						ServerId: body.ServerID,
					},
					From: int32(body.From),
					To:   int32(body.To),
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) AddRoleMember(body *types.ChangeRoleMemberParams) {
	serversPID := se.GetAllServerInstances(body.ServerID)

	for _, serverPID := range serversPID {
		message := &message.WSMessage{
			Content: &message.WSMessage_AddRoleMember{
				AddRoleMember: &message.AddRoleMember{
					UserId: body.UserID,
					Role: &message.Role{
						Id:       body.RoleID,
						ServerId: body.ServerID,
					},
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) RemoveRoleMember(body *types.ChangeRoleMemberParams) {
	serversPID := se.GetAllServerInstances(body.ServerID)

	for _, serverPID := range serversPID {
		message := &message.WSMessage{
			Content: &message.WSMessage_RemoveRoleMember{
				RemoveRoleMember: &message.RemoveRoleMember{
					UserId: body.UserID,
					Role: &message.Role{
						Id:       body.RoleID,
						ServerId: body.ServerID,
					},
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) SendFriendRequest(friendshipID, receiverID string, sender *db.User) {
	userPID := se.GetUser(receiverID)

	message := &message.WSMessage{
		Content: &message.WSMessage_FriendRequest{
			FriendRequest: &message.FriendRequest{
				FriendshipId: friendshipID,
				Sender: &message.User{
					Id:          sender.ID,
					DisplayName: sender.DisplayName,
					Avatar:      sender.Avatar.String,
					Banner:      sender.Banner.String,
					AboutMe:     sender.AboutMe,
				},
				Accepted: false,
			},
		},
	}

	se.cluster.Engine().Send(userPID, message)
}

func (se *service) AcceptFriendRequest(friendshipID, senderID, receiverID, channelID string) {
	senderPID := se.GetUser(senderID)
	receiverPID := se.GetUser(receiverID)

	message := &message.WSMessage{
		Content: &message.WSMessage_AcceptFriendRequest{
			AcceptFriendRequest: &message.AcceptFriendRequest{
				FriendshipId: friendshipID,
				ChannelId:    channelID,
			},
		},
	}

	se.StartDMChannel(channelID, []string{senderID, receiverID})

	se.cluster.Engine().Send(senderPID, message)
	se.cluster.Engine().Send(receiverPID, message)
}

func (se *service) RemoveFriend(friendshipID, senderID, receiverID, channelID string) {
	senderPID := se.GetUser(senderID)
	receiverPID := se.GetUser(receiverID)

	message := &message.WSMessage{
		Content: &message.WSMessage_RemoveFriend{
			RemoveFriend: &message.RemoveFriend{
				FriendshipId: friendshipID,
			},
		},
	}

	channelPIDs := se.GetAllChannelInstances("global", channelID)
	for _, channelPID := range channelPIDs {
		se.cluster.Engine().Poison(channelPID)
	}

	se.cluster.Engine().Send(senderPID, message)
	se.cluster.Engine().Send(receiverPID, message)
}

func (se *service) SendUserStatusMessage(userPID *actor.PID, status *message.ChangeStatus) {
	servers := se.GetAllServerInstances(status.ServerId)
	for _, serverPID := range servers {
		se.cluster.Engine().SendWithSender(serverPID, status, userPID)
	}
}

func (se *service) NotifyAccountDeletion(userID string, serverIDs []string) {
	userPID := se.GetUser(userID)

	for _, serverID := range serverIDs {
		serverPIDs := se.GetAllServerInstances(serverID)
		for _, serverPID := range serverPIDs {
			se.cluster.Engine().Send(serverPID, &message.AccountDeletion{
				UserId:   userID,
				ServerId: serverID,
			})
		}
	}

	se.cluster.Engine().Send(userPID, &message.AccountDeletion{
		UserId: userID,
	})
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

func (se *service) GetActiveFriends(userID string) []string {
	var friendIDs []string

	userPID := se.GetUser(userID)
	response := se.cluster.Engine().Request(userPID, &message.GetFriends{}, 10*time.Second)
	result, err := response.Result()
	if err == nil {
		friendIDs = append(friendIDs, result.(*message.GetFriends).FriendIds...)
	}

	return friendIDs
}

func (se *service) AvatarServerChange(serverID string, bannerURL, avatarURL *string) {
	serverPIDs := se.GetAllServerInstances(serverID)

	for _, serverPID := range serverPIDs {
		message := &message.WSMessage{
			Content: &message.WSMessage_AvatarServerChange{
				AvatarServerChange: &message.AvatarServerChange{
					ServerId:  serverID,
					AvatarUrl: avatarURL,
					BannerUrl: bannerURL,
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) ProfileServerChange(serverID string, body *types.UpdateServerProfileParams) {
	serverPIDs := se.GetAllServerInstances(serverID)

	for _, serverPID := range serverPIDs {
		message := &message.WSMessage{
			Content: &message.WSMessage_ProfileServerChange{
				ProfileServerChange: &message.ProfileServerChange{
					ServerId:    serverID,
					Name:        body.Name,
					Description: body.Description,
					Public:      body.Public,
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) EditChannel(channelID string, body *types.EditChannelParams) {
	serverPIDs := se.GetAllServerInstances(body.ServerID)

	for _, serverPID := range serverPIDs {
		message := &message.WSMessage{
			Content: &message.WSMessage_EditChannel{
				EditChannel: &message.EditChannel{
					Channel: &message.Channel{
						Id:          channelID,
						ServerId:    body.ServerID,
						Name:        body.Name,
						Description: body.Description,
						Users:       body.Users,
						Roles:       body.Roles,
					},
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) EditCategory(categoryID string, body *types.EditCategoryParams) {
	serverPIDs := se.GetAllServerInstances(body.ServerID)

	for _, serverPID := range serverPIDs {
		message := &message.WSMessage{
			Content: &message.WSMessage_EditCategory{
				EditCategory: &message.EditCategory{
					Category: &message.Category{
						Id:       categoryID,
						ServerId: body.ServerID,
						Name:     body.Name,
						Users:    body.Users,
						Roles:    body.Roles,
					},
				},
			},
		}

		se.cluster.Engine().Send(serverPID, message)
	}
}

func (se *service) BanUser(serverID string, body *types.BanUserParams) {
	serverPIDs := se.GetAllServerInstances(serverID)

	for _, serverPID := range serverPIDs {
		se.cluster.Engine().Send(serverPID, &message.BanUser{
			ServerId: serverID,
			UserId:   body.UserID,
			Reason:   body.Reason,
			Duration: timestamppb.New(body.Duration),
		})
	}
}

func (se *service) KickUser(serverID string, body *types.KickUserParams) {
	serverPIDs := se.GetAllServerInstances(serverID)

	for _, serverPID := range serverPIDs {
		se.cluster.Engine().Send(serverPID, &message.KickUser{
			ServerId: serverID,
			UserId:   body.UserID,
			Reason:   body.Reason,
		})
	}
}
