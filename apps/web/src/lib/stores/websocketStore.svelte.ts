import { WSMessageSchema } from '$lib/gen/types_pb';
import type { Category, Channel, ChannelTypes, Friend, Member, Message, Role } from '$lib/types/types';
import { fromBinary } from '@bufbuild/protobuf';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { channelStore } from './channelStore.svelte';
import { page } from '$app/state';
import { serverStore } from './serverStore.svelte';
import { userStore } from './userStore.svelte';
import { categoryStore } from './categoryStore.svelte';
import { goto } from '$app/navigation';

export class WebsocketStore {
  wsConn = $state<WebSocket>();

  init(userID: string) {
    const ws = new WebSocket(`ws://localhost:8080/api/protected/ws/${userID}`);
    if (!ws) return;

    this.wsConn = ws;

    ws.onopen = () => {
      window.setInterval(() => {
        ws.send('heartbeat');
      }, 10 * 1000);
    };

    ws.onmessage = async (event) => {
      if (event.data === 'heartbeat') return;

      const arrayBuffer = await event.data.arrayBuffer();
      const uint8Array = new Uint8Array(arrayBuffer);
      const wsMess = fromBinary(WSMessageSchema, uint8Array, {
        readUnknownFields: false
      });

      switch (wsMess.content.case) {
        case 'userChangeStatus':
          {
            if (!wsMess.content.value) return;
            if (wsMess.content.value.user?.id === userStore.user?.id) return;
            if (page.params.server_id !== wsMess.content.value.serverId) return;
            const value = wsMess.content.value;

            if (value.status === 'offline') {
              serverStore.setMemberOffline(value.serverId, value.user!.id);
            } else {
              if (value.type === 'connect')
                serverStore.setMemberOnline(value.serverId, value.user!.id, value.status);
              if (value.type === 'join') {
                const member: Member = {
                  id: value.user!.id,
                  display_name: value.user!.displayName,
                  avatar: value.user!.avatar,
                  status: value.status,
                  roles: []
                };
                serverStore.addMember(value.serverId, member);
              }
            }
          }
          break;
        case 'newChatMessage':
          {
            if (!wsMess.content.value.message) return;
            if (wsMess.content.value.message.channelId !== page.params.channel_id) return;
            const message = wsMess.content.value.message;
            const contentStr = new TextDecoder().decode(message.content);
            const attachments = new TextDecoder().decode(message.attachments);

            const newMessage: Message = {
              id: message.id,
              author: {
                id: message.author!.id,
                avatar: message.author!.avatar,
                display_name: message.author!.displayName
              },
              server_id: message.serverId,
              channel_id: message.channelId,
              content: JSON.parse(contentStr),
              everyone: message.everyone,
              mentions_users: message.mentionsUsers,
              mentions_channels: message.mentionsChannels,
              attachments: attachments.length > 0 ? JSON.parse(attachments) : [],
              updated_at: timestampDate(message.createdAt!).toISOString(),
              created_at: timestampDate(message.createdAt!).toISOString()
            };

            channelStore.addMessage(newMessage);
          }
          break;
        case 'deleteChatMessage':
          {
            if (!wsMess.content.value.message) return;
            if (wsMess.content.value.message.channelId !== page.params.channel_id) return;
            const message = wsMess.content.value.message;
            channelStore.deleteMessage(message.id);
          }
          break;
        case 'editChatMessage':
          {
            if (!wsMess.content.value.message) return;
            if (wsMess.content.value.message.channelId !== page.params.channel_id) return;

            const message = wsMess.content.value.message;
            const contentStr = new TextDecoder().decode(message.content);

            const editMessage: Partial<Message> = {
              id: message.id,
              everyone: message.everyone,
              mentions_users: message.mentionsUsers,
              mentions_channels: message.mentionsChannels,
              content: JSON.parse(contentStr),
              updated_at: timestampDate(message.updatedAt!).toISOString()
            };

            channelStore.editMessage(editMessage);
          }
          break;
        case 'startCategory':
          {
            if (!wsMess.content.value.category) return;
            const category = wsMess.content.value.category;

            const newCategory: Category = {
              id: category.id,
              server_id: category.serverId,
              name: category.name,
              position: category.position,
              users: category.users,
              roles: category.roles,
              e2ee: category.e2ee,
              channels: {}
            };

            categoryStore.addCategory(newCategory);
          }
          break;
        case 'startChannel':
          {
            if (!wsMess.content.value.channel) return;
            const channel = wsMess.content.value.channel;

            const newChannel: Channel = {
              id: channel.id,
              server_id: channel.serverId,
              category_id: channel.categoryId,
              name: channel.name,
              description: channel.description,
              users: channel.users,
              roles: channel.roles,
              type: channel.type as ChannelTypes,
              position: channel.position,
              unread: false
            };

            channelStore.addChannel(newChannel);
          }
          break;
        case 'killCategory':
          {
            if (!wsMess.content.value) return;
            const value = wsMess.content.value;
            const channels = channelStore.getCategoryChannels(value.serverId, value.categoryId);

            if (channels.find((chan) => chan.id === page.params.channel_id)) {
              const firstChan = channelStore.getFirstChannel(value.serverId);
              if (firstChan) goto(`/servers/${value.serverId}/channels/${firstChan}`);
            }

            categoryStore.deleteCategory(value.serverId, value.categoryId);
          }
          break;
        case 'killChannel':
          {
            if (!wsMess.content.value.channel) return;
            const channel = wsMess.content.value.channel;

            if (channel.id === page.params.channel_id) {
              const firstChan = channelStore.getFirstChannel(channel.serverId);
              if (firstChan) goto(`/servers/${channel.serverId}/channels/${firstChan}`);
            }

            channelStore.deleteChannel(channel.serverId, channel.categoryId, channel.id);
          }
          break;
        case 'createOrEditRole':
          {
            if (!wsMess.content.value.role) return;
            const role = wsMess.content.value.role;
            const newRole: Role = {
              id: role.id,
              members: [],
              name: role.name,
              color: role.color,
              abilities: role.abilities,
              position: role.position
            };

            const existingRole = serverStore.getRole(role.serverId, role.id);
            if (existingRole) {
              serverStore.editRole(role.serverId, newRole);
            } else {
              serverStore.addRole(role.serverId, newRole);
            }
          }
          break;
        case 'addRoleMember':
          {
            if (!wsMess.content.value.userId || !wsMess.content.value.role) return;
            const userId = wsMess.content.value.userId;
            const role = wsMess.content.value.role;

            const member = serverStore.getMember(role.serverId, userId);
            if (member) member.roles.push(role.id);

            if (userId === userStore.user!.id) {
              serverStore.servers[role.serverId].user_roles.push(role.id);
            }
          }
          break;
        case 'removeRoleMember':
          {
            if (!wsMess.content.value.userId || !wsMess.content.value.role) return;
            const userId = wsMess.content.value.userId;
            const role = wsMess.content.value.role;

            const member = serverStore.getMember(role.serverId, userId);
            if (member) {
              member.roles = member.roles.filter((roleID) => roleID !== role.id);
            }

            if (userId === userStore.user!.id) {
              serverStore.servers[role.serverId].user_roles = serverStore.servers[
                role.serverId
              ].user_roles.filter((roleID) => roleID !== role.id);
            }
          }
          break;
        case 'removeRole':
          {
            if (!wsMess.content.value.role) return;
            const role = wsMess.content.value.role;
            serverStore.deleteRole(role.serverId, role.id);
          }
          break;
        case 'moveRole':
          {
            if (!wsMess.content.value) return;
            const movedRoleID = wsMess.content.value.movedRole?.id;
            const targetRoleID = wsMess.content.value.targetRole?.id;
            const serverID = wsMess.content.value.movedRole?.serverId;
            if (!movedRoleID || !targetRoleID || !serverID) return;

            const fromPos = wsMess.content.value.from;
            const toPos = wsMess.content.value.to;

            const roles = serverStore.getRoles(serverID);
            const movedRole = serverStore.getRole(serverID, movedRoleID);
            const targetRole = serverStore.getRole(serverID, targetRoleID);

            if (movedRole && targetRole) {
              movedRole.position = toPos;
              targetRole.position = fromPos;
            }

            roles.sort((a, b) => a.position - b.position);
          }
          break;
        case 'friendRequest':
          {
            if (!wsMess.content.value) return;
            const sender = wsMess.content.value.sender
            const friendshipID = wsMess.content.value.friendshipId
            const accepted = wsMess.content.value.accepted
            if (!sender) return;

            const aboutMeStr = new TextDecoder().decode(sender.aboutMe);

            const friend: Friend = {
              friendship_id: friendshipID,
              friendship_sender_id: sender.id,
              id: sender.id,
              display_name: sender.displayName,
              about_me: JSON.parse(aboutMeStr),
              banner: sender.banner,
              avatar: sender.avatar,
              accepted: accepted
            }

            userStore.friends.push(friend)
          }
          break;
        case 'acceptFriendRequest':
          {
            if (!wsMess.content.value) return;
            const friendshipID = wsMess.content.value.friendshipId
            const channelID = wsMess.content.value.channelId

            userStore.acceptFriend(friendshipID, channelID)
          }
          break;
        case 'removeFriend':
          {
            if (!wsMess.content.value) return;
            const friendshipID = wsMess.content.value.friendshipId

            userStore.removeFriend(friendshipID)
          }
          break;

      }
    };
  }
}

export const ws = new WebsocketStore();
