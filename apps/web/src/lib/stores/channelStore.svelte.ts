import type { Channel, Message } from '$lib/types/types';
import { logErr } from 'utils/print';
import { backend } from './backendStore.svelte';
import { categoryStore } from './categoryStore.svelte';
import { serverStore } from './serverStore.svelte';

const AVG_MESSAGE_HEIGHT = 100;

class ChannelStore {
	messageCache = $state<
		Record<
			string,
			{
				beforeMessageID: string;
				afterMessageID: string;
				messages: Message[];
				hasReachedEnd: boolean;
				offsetHeight: number;
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

	messageIsRecent(channelID: string, messageID: string) {
		if (!this.messageCache[channelID]) return;

		const messageIdx = this.messageCache[channelID].messages.findIndex(
			(message) => message.id === messageID
		);
		if (messageIdx < 0) return false;

		const message = this.messageCache[channelID].messages[messageIdx];
		const nextMessage = this.messageCache[channelID].messages[messageIdx - 1];

		if (!nextMessage || !message) return false;
		if (nextMessage.author.id !== message.author.id) return false;

		const messageDate = new Date(message.created_at);
		const nextMessageDate = new Date(nextMessage.created_at);

		const diff = (nextMessageDate.getTime() - messageDate.getTime()) / 1000;

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
		const channel = this.getChannel(channelOpts.server_id!, channelID);
		if (!channel) return;

		if (channelOpts.name) channel.name = channelOpts.name;
		if (channelOpts.description) channel.description = channelOpts.description;
		if (channelOpts.users) channel.users = channelOpts.users;
		if (channelOpts.roles) channel.roles = channelOpts.roles;
	}

	addMessage(channelID: string, message: Message) {
		if (!this.messageCache[channelID]) {
			this.messageCache[channelID] = {
				messages: [],
				lastMessageID: '',
				hasReachedEnd: false,
				offsetHeight: 0
			};
		}
		this.messageCache[channelID].messages.unshift(message);
	}

	editMessage(channelID: string, message: Partial<Message>) {
		if (!this.messageCache[channelID]) return;
		const index = this.messageCache[channelID].messages.findIndex((m) => m.id === message.id);
		if (index !== -1) {
			this.messageCache[channelID].messages[index] = {
				...this.messageCache[channelID].messages[index],
				...message
			};
		}
	}

	deleteMessage(channelID: string, messageID: string) {
		if (!this.messageCache[channelID]) return;
		this.messageCache[channelID].messages = this.messageCache[channelID].messages.filter(
			(message) => message.id !== messageID
		);
	}

	deleteChannel(serverID: string, categoryID: string, channelID: string) {
		const category = categoryStore.getCategory(serverID, categoryID);
		if (!category) return null;

		delete category.channels[channelID];
	}

	private initializeMessageCache(channelID: string) {
		this.messageCache[channelID] = {
			messages: [],
			beforeMessageID: '',
			afterMessageID: '',
			hasReachedEnd: false,
			offsetHeight: 0
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
			cache.offsetHeight += removedCount * AVG_MESSAGE_HEIGHT;
		} else {
			cache.messages = cache.messages.slice(0, TRIM_AMOUNT);
			cache.offsetHeight -= removedCount * AVG_MESSAGE_HEIGHT;
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

	async loadMoreMessages(
		serverID: string,
		channelID: string,
		direction: 'before' | 'after'
	): Promise<boolean> {
		if (!this.messageCache[channelID]) {
			this.initializeMessageCache(channelID);
		}

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

		return true;
	}

	async ensureMessagesLoaded(serverID: string, channelID: string): Promise<void> {
		if (!this.messageCache[channelID] || this.messageCache[channelID].messages.length === 0) {
			await this.loadMoreMessages(serverID, channelID, 'before');
		}
	}

	clearChannelCache(channelID: string) {
		delete this.messageCache[channelID];
	}
}

export const channelStore = new ChannelStore();
