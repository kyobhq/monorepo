import type { Channel, Message } from '$lib/types/types';
import { categoryStore } from './categoryStore.svelte';
import { serverStore } from './serverStore.svelte';

class ChannelStore {
  currentChannel = $state<Channel | null>(null);
  messages = $state<Message[]>([]);

  getFirstChannel(serverID: string) {
    const server = serverStore.getServer(serverID);
    const firstCategory = Object.values(server?.categories || {}).find(
      (category) => category.position === 0
    );
    const firstChannel = Object.values(firstCategory?.channels || {}).find(
      (chan) => chan.position === 0
    );
    return firstChannel?.id || null;
  }

  messageIsRecent(messageID: string) {
    const messageIdx = this.messages.findIndex((message) => message.id === messageID);
    if (messageIdx < 0) return false;

    const message = this.messages[messageIdx];
    const nextMessage = this.messages[messageIdx - 1];

    if (!nextMessage || !message) return false;
    if (nextMessage.author.id !== message.author.id) return false;

    const messageDate = new Date(message.created_at);
    const nextMessageDate = new Date(nextMessage.created_at);

    const diff = (nextMessageDate.getTime() - messageDate.getTime()) / 1000

    return diff < 30;
  }

  getCategoryChannels(serverID: string, categoryID: string): Channel[] {
    const server = serverStore.getServer(serverID);
    if (!server?.categories) return [];
    const category = categoryStore.getCategory(serverID, categoryID);

    return Object.values(category?.channels || []);
  }

  getChannel(serverID: string, channelID: string) {
    const server = serverStore.getServer(serverID);
    if (!server?.categories) return null;

    for (const category of Object.values(server.categories)) {
      const channel = category.channels?.[channelID];
      if (channel) return channel;
    }

    return null;
  }

  getChannelsLastPositionInCategory(serverID: string, categoryID: string) {
    const category = categoryStore.getCategory(serverID, categoryID);
    return Object.values(category?.channels || {}).length;
  }

  addChannel(channel: Channel) {
    const category = categoryStore.getCategory(channel.server_id, channel.category_id);
    if (!category) return null;

    category.channels[channel.id] = channel;
  }

  editChannel(channelID: string, channelOpts: Partial<Channel>) {
    const channel = this.getChannel(channelOpts.server_id!, channelID)
    if (!channel) return;

    if (channelOpts.name) channel.name = channelOpts.name
    if (channelOpts.description) channel.description = channelOpts.description
    if (channelOpts.users) channel.users = channelOpts.users
    if (channelOpts.roles) channel.roles = channelOpts.roles
  }

  addMessage(message: Message) {
    this.messages.unshift(message);
  }

  editMessage(message: Partial<Message>) {
    const index = this.messages.findIndex((m) => m.id === message.id);
    if (index !== -1) {
      this.messages[index] = {
        ...this.messages[index],
        ...message
      };
    }
  }

  deleteMessage(messageID: string) {
    this.messages = this.messages.filter((message) => message.id !== messageID);
  }

  deleteChannel(serverID: string, categoryID: string, channelID: string) {
    const category = categoryStore.getCategory(serverID, categoryID);
    if (!category) return null;

    delete category.channels[channelID];
  }
}

export const channelStore = new ChannelStore();
