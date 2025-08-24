import type { Message, Member } from '$lib/types/types';

export class MessageStore {
	editMessage = $state<Message | null>(null);
	messageAuthors = $state<Record<string, Member>>({});

	stopEditing(): void {
		this.editMessage = null;
	}

	cacheAuthor(author: Member): void {
		if (author.id) {
			this.messageAuthors[author.id] = author;
		}
	}

	cacheAuthors(authors: Member[]): void {
		for (const author of authors) {
			this.cacheAuthor(author);
		}
	}

	getAuthor(authorId: string): Member | undefined {
		return this.messageAuthors[authorId];
	}

	updateAuthor(authorId: string, updates: Partial<Member>): void {
		const author = this.messageAuthors[authorId];
		if (author) {
			this.messageAuthors[authorId] = { ...author, ...updates };
		}
	}

	cacheMessageAuthors(messages: Message[]): void {
		for (const message of messages) {
			this.cacheAuthor(message.author);
		}
	}

	clearAuthorCache(): void {
		this.messageAuthors = {};
	}

	removeAuthor(authorId: string): void {
		delete this.messageAuthors[authorId];
	}

	getAuthorRoles(authorId: string): string[] {
		return this.messageAuthors[authorId]?.roles || [];
	}
}

export const messageStore = new MessageStore();
