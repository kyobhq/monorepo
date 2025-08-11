import { WSMessageSchema } from '$lib/gen/types_pb';
import type { Channel, ChannelTypes, Member, Message } from '$lib/types/types';
import { fromBinary } from '@bufbuild/protobuf';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { channelStore } from './channelStore.svelte';
import { page } from '$app/state';
import { serverStore } from './serverStore.svelte';
import { userStore } from './userStore.svelte';

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
              unread: false,
            }

            channelStore.addChannel(newChannel)
          }
          break;
        case 'killChannel':
          {
            if (!wsMess.content.value.channel) return;
            const channel = wsMess.content.value.channel;

            channelStore.deleteChannel(channel.serverId, channel.categoryId, channel.id);
          }
          break;
      }
    };
  }
}

export const ws = new WebsocketStore();
