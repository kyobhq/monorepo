import type { Channel, Message } from '$lib/types/types';
import { logErr } from 'utils/print';
import { backend } from './backendStore.svelte';
import { categoryStore } from './categoryStore.svelte';
import { serverStore } from './serverStore.svelte';
import { messageStore } from './messageStore.svelte';
import { page } from '$app/state';
import { userStore } from './userStore.svelte';

const AVG_MESSAGE_HEIGHT = 100;
const SCROLL_LIMIT = 3000;

class ChannelStore {
  messageCache = $state<
    Record<
      string,
      {
        beforeMessageID: string;
        afterMessageID: string;
        messages: Message[];
        hasReachedEnd: boolean;
        scrollHeight: number;
        scrollY: number;
      }
    >
  >({});

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

  messageIsRecent(channelID: string, messageID: string): boolean {
    const cache = this.messageCache[channelID];
    if (!cache) return false;

    const messageIdx = cache.messages.findIndex((m) => m.id === messageID);
    if (messageIdx < 0) return false;

    const [message, prevMessage] = [cache.messages[messageIdx], cache.messages[messageIdx + 1]];
    if (!message || !prevMessage || prevMessage.author.id !== message.author.id) {
      return false;
    }

    const timeDiff =
      (new Date(message.created_at).getTime() - new Date(prevMessage.created_at).getTime()) / 1000;
    return timeDiff < 30;
  }

  getCategoryChannels(serverID: string, categoryID: string): Channel[] {
    const category = categoryStore.getCategory(serverID, categoryID);
    return Object.values(category?.channels || {});
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

  getChannelsLastPositionInCategory(serverID: string, categoryID: string): number {
    return this.getCategoryChannels(serverID, categoryID).length;
  }

  addChannel(channel: Channel): void {
    const category = categoryStore.getCategory(channel.server_id, channel.category_id);
    if (category) {
      category.channels[channel.id] = channel;
    }
  }

  editChannel(channelID: string, channelOpts: Partial<Channel>): void {
    const channel = this.getChannel(channelOpts.server_id!, channelID);
    if (!channel) return;

    if (channelOpts.name) channel.name = channelOpts.name;
    if (channelOpts.description) channel.description = channelOpts.description;
    if (channelOpts.users) channel.users = channelOpts.users;
    if (channelOpts.roles) channel.roles = channelOpts.roles;
  }

  addMessage(serverID: string, channelID: string, message: Message): void {
    const channel = this.getChannel(serverID, channelID);
    if (!channel) return;

    const isCurrentChannel = page.params.channel_id === channelID;

    if (isCurrentChannel && !this.messageCache[channelID]) {
      this.initializeMessageCache(channelID);
    }

    if (this.messageCache[channelID]) {
      this.messageCache[channelID].messages.unshift(message);
      messageStore.cacheAuthor(message.author);
    }

    channel.last_message_sent = message.id;

    if (isCurrentChannel) channel.last_message_read = message.id;
    if (message.mentions_users.includes(userStore.user!.id)) {
      channel.last_mentions = channel.last_mentions || [];
      channel.last_mentions.push(message.id);
    }
  }

  addMessageDM(channelID: string, message: Message): void {
    const friend = userStore.getFriendByChannelID(channelID);
    if (!friend) return;

    const isCurrentChannel = page.params.channel_id === channelID;

    if (isCurrentChannel && !this.messageCache[channelID]) {
      this.initializeMessageCache(channelID);
    }

    if (this.messageCache[channelID]) {
      this.messageCache[channelID].messages.unshift(message);
      messageStore.cacheAuthor(message.author);
    }

    friend.last_message_sent = message.id;
    if (isCurrentChannel) friend.last_message_read = message.id;
  }

  setLastMessageRead(serverID: string, channelID: string): void {
    const cache = this.messageCache[channelID];
    if (!cache || cache.scrollY > SCROLL_LIMIT) return;

    const latestMessageId = cache.messages[0]?.id;
    if (!latestMessageId) return;

    if (serverID === 'global') {
      const friend = userStore.getFriendByChannelID(channelID);
      if (friend) {
        friend.last_message_read = latestMessageId;
      }
      return;
    }

    const channel = this.getChannel(serverID, channelID);
    if (channel) {
      channel.last_message_read = latestMessageId;
      channel.last_mentions = [];
    }
  }

  setLastMessageSent(serverID: string, channelID: string): void {
    const cache = this.messageCache[channelID];
    const latestMessageId = cache?.messages[0]?.id;
    if (!latestMessageId) return;

    if (serverID === 'global') {
      const friend = userStore.getFriendByChannelID(channelID);
      if (friend) {
        friend.last_message_sent = latestMessageId;
      }
      return;
    }

    const channel = this.getChannel(serverID, channelID);
    if (channel && !channel.last_message_sent) {
      channel.last_message_sent = latestMessageId;
    }
  }

  editMessage(channelID: string, message: Partial<Message>): void {
    const cache = this.messageCache[channelID];
    if (!cache) return;

    const index = cache.messages.findIndex((m) => m.id === message.id);
    if (index !== -1) {
      cache.messages[index] = { ...cache.messages[index], ...message };
    }
  }

  deleteMessage(channelID: string, messageID: string): void {
    const cache = this.messageCache[channelID];
    if (cache) {
      cache.messages = cache.messages.filter((message) => message.id !== messageID);
    }
  }

  deleteChannel(serverID: string, categoryID: string, channelID: string): void {
    const category = categoryStore.getCategory(serverID, categoryID);
    if (category) {
      delete category.channels[channelID];
    }
  }

  private initializeMessageCache(channelID: string): void {
    this.messageCache[channelID] = {
      messages: [],
      beforeMessageID: '',
      afterMessageID: '',
      hasReachedEnd: false,
      scrollHeight: 0,
      scrollY: 0
    };
  }

  private async fetchMessages(serverID: string, channelID: string, direction: 'before' | 'after') {
    const cache = this.messageCache[channelID];
    const messageID = direction === 'before' ? cache.beforeMessageID : cache.afterMessageID;

    return await backend.getMessages({
      serverID,
      channelID,
      ...(direction === 'before' ? { beforeMessageID: messageID } : { afterMessageID: messageID })
    });
  }

  private trimCacheIfNeeded(channelID: string, direction: 'before' | 'after') {
    const cache = this.messageCache[channelID];
    const MAX_CACHE_SIZE = 100;
    const TRIM_AMOUNT = 50;

    if (cache.messages.length < MAX_CACHE_SIZE) return;

    const removedCount = cache.messages.length - TRIM_AMOUNT;
    if (direction === 'before') {
      cache.messages = cache.messages.slice(TRIM_AMOUNT);
      cache.scrollHeight += removedCount * AVG_MESSAGE_HEIGHT;
    } else {
      cache.messages = cache.messages.slice(0, TRIM_AMOUNT);
      cache.scrollHeight -= removedCount * AVG_MESSAGE_HEIGHT;
    }
  }

  private addMessagesToCache(
    channelID: string,
    messages: Message[],
    direction: 'before' | 'after'
  ) {
    const cache = this.messageCache[channelID];
    if (direction === 'before') {
      cache.messages.push(...messages);
    } else {
      cache.messages.unshift(...messages);
    }

    cache.beforeMessageID = cache.messages[cache.messages.length - 1].id;
    cache.afterMessageID = cache.messages[0].id;
  }

  private cacheMessageAuthors(messages: Message[]) {
    messageStore.cacheMessageAuthors(messages);
  }

  async loadMoreMessages(
    serverID: string,
    channelID: string,
    direction: 'before' | 'after'
  ): Promise<boolean> {
    if (!this.messageCache[channelID]) this.initializeMessageCache(channelID);

    const cache = this.messageCache[channelID];

    if (cache.hasReachedEnd) return false;

    const res = await this.fetchMessages(serverID, channelID, direction);
    if (res.isErr()) {
      logErr(res.error);
      return false;
    }

    const messages = res.value;
    if (!messages?.length) {
      cache.hasReachedEnd = true;
      return false;
    }

    this.trimCacheIfNeeded(channelID, direction);
    this.addMessagesToCache(channelID, messages, direction);
    this.cacheMessageAuthors(messages);

    return true;
  }

  async ensureMessagesLoaded(serverID: string, channelID: string): Promise<void> {
    if (!this.messageCache[channelID] || this.messageCache[channelID].messages.length === 0) {
      await this.loadMoreMessages(serverID, channelID, 'before');
    }
  }

  clearChannelCache(channelID: string): void {
    delete this.messageCache[channelID];
  }
}

export const channelStore = new ChannelStore();
