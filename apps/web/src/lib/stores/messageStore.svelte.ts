import type { Message, Member } from '$lib/types/types';
import { serverStore } from './serverStore.svelte';

export class MessageStore {
	editMessage = $state<Message | null>(null);
	messageAuthors = $state<Record<string, Member>>({});

	stopEditing() {
		this.editMessage = null;
	}

	cacheAuthor(author: Member) {
		if (!author.id) return;
		this.messageAuthors[author.id] = author;
	}

	cacheAuthors(authors: Member[]) {
		for (const author of authors) {
			this.cacheAuthor(author);
		}
	}

	getAuthor(authorId: string): Member | undefined {
		return this.messageAuthors[authorId];
	}

	updateAuthor(authorId: string, updates: Partial<Member>) {
		if (this.messageAuthors[authorId]) {
			this.messageAuthors[authorId] = { ...this.messageAuthors[authorId], ...updates };
		}
	}

	cacheMessageAuthors(messages: Message[]) {
		for (const message of messages) {
			this.cacheAuthor(message.author);
		}
	}

	clearAuthorCache() {
		this.messageAuthors = {};
	}

	removeAuthor(authorId: string) {
		delete this.messageAuthors[authorId];
	}

	getAuthorRoles(authorId: string): string[] {
		const cachedAuthor = this.messageAuthors[authorId];

		if (cachedAuthor?.roles?.length) {
			return cachedAuthor.roles;
		}

		return [];
	}
}

export const messageStore = new MessageStore();
