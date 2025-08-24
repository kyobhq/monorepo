import type { EditServerType } from '$lib/types/schemas';
import type { Invite, Member, Role, Server } from '$lib/types/types';
import { createId } from '@paralleldrive/cuid2';
import { backend } from './backendStore.svelte';
import { logErr } from 'utils/print';
import { userStore } from './userStore.svelte';
import type { Abilities } from '$lib/constants/permissions';
import { page } from '$app/state';
import { messageStore } from './messageStore.svelte';
import { channelStore } from './channelStore.svelte';

interface CacheEntry {
	cached: boolean;
	lastAccessed: number;
	timeoutId?: NodeJS.Timeout;
}

export class ServerStore {
	servers = $state<Record<string, Server>>({});
	memberCount = $state<number>(0);
	cached = $state<Record<string, CacheEntry>>({});
	abilities = $derived.by(() => {
		const allAbilities: Record<string, Abilities[]> = {};

		for (const server of Object.values(this.servers)) {
			const user_roles = server?.user_roles || [];
			allAbilities[server.id] = [];

			for (const user_role of user_roles) {
				const role = this.getRole(server.id, user_role);
				if (role?.abilities) allAbilities[server.id].push(...role.abilities);
			}

			if (userStore.user!.id === server.owner_id) allAbilities[server.id].push('OWNER');
		}

		return allAbilities;
	});

	safeServerOperation<T>(serverID: string, operation: (server: Server) => T, fallback: T): T {
		const server = this.servers[serverID];
		if (!server) {
			return fallback;
		}

		try {
			return operation(server);
		} catch (error) {
			console.warn(`Server operation failed for ${serverID}:`, error);
			return fallback;
		}
	}

	getServer(id: string) {
		return this.servers[id];
	}

	addMember(serverID: string, member: Member) {
		if (this.getMember(serverID, member.id!)) return;
		if (page.params.server_id === serverID) serverStore.memberCount += 1;

		this.servers[serverID].members.push(member);
	}

	addMembers(serverID: string, members: Member[]) {
		this.servers[serverID].members.push(...members);
	}

	resetMembersList(serverID: string) {
		this.servers[serverID].members = this.servers[serverID].members.slice(0, 50);
	}

	setMemberOnline(serverID: string, memberID: string, status: string) {
		if (page.params.server_id !== serverID || serverID === '') return;

		const member = this.servers[serverID].members?.find((m) => m.id === memberID);
		if (member) member.status = status;
	}

	setMemberOffline(serverID: string, memberID: string) {
		if (page.params.server_id !== serverID || serverID === '') return;

		const member = this.servers[serverID].members?.find((m) => m.id === memberID);
		if (member) member.status = 'offline';
	}

	deleteMember(serverID: string, userID: string) {
		this.servers[serverID].members = this.servers[serverID].members.filter((m) => m.id !== userID);
	}

	getMember(serverID: string, userID: string) {
		return this.servers[serverID].members.find((m) => m.id === userID);
	}

	getMemberRoles(serverID: string, userID: string) {
		const cachedAuthor = messageStore.getAuthor(userID);
		if (cachedAuthor) {
			return cachedAuthor.roles;
		}

		const member = this.servers[serverID].members.find((m) => m.id === userID);
		if (member) {
			return member.roles;
		}

		return [];
	}

	setMembers(serverID: string, members: Member[]) {
		this.servers[serverID].members = members;
	}

	getMembers(serverID: string) {
		return this.servers[serverID].members;
	}

	setRoles(serverID: string, roles: Role[]) {
		this.servers[serverID].roles = roles?.sort((a, b) => a.position - b.position) || [];
	}

	setUserRoles(serverID: string, roles: string[]) {
		this.servers[serverID].user_roles = roles;
	}

	getUserRoles(serverID: string) {
		return this.servers[serverID].user_roles || [];
	}

	getRoles(serverID: string) {
		return this.servers[serverID].roles;
	}

	getRole(serverID: string, roleID: string): Role | undefined {
		const roles = this.servers[serverID].roles;

		if (roleID === 'default') {
			return {
				id: roleID,
				color: '',
				name: 'Members',
				position: roles.length,
				abilities: [],
				members: []
			};
		} else if (roleID === 'offline') {
			return {
				id: roleID,
				color: '',
				name: 'Offline',
				position: roles.length,
				abilities: [],
				members: []
			};
		}

		return roles.find((role) => role.id === roleID);
	}

	addRole(serverID: string, role: Role) {
		this.servers[serverID].roles.push(role);
	}

	editRole(serverID: string, role: Role) {
		const roles = this.servers[serverID].roles;
		const roleIdx = roles.findIndex((r) => r.id === role.id);
		if (roleIdx !== -1) {
			roles[roleIdx] = { ...roles[roleIdx], ...role };
			this.servers[serverID].roles = [...roles];
		}
	}

	deleteRole(serverID: string, roleID: string) {
		this.servers[serverID].roles = this.servers[serverID].roles.filter(
			(role) => role.id !== roleID
		);
	}

	async createTemplateRole(serverID: string) {
		const roles = this.getRoles(serverID);

		const newRole: Role = {
			id: createId(),
			name: 'new role',
			color: '#ADADB8',
			abilities: [],
			position: roles.length,
			members: []
		};
		this.servers[serverID].roles.push(newRole);

		const res = await backend.createOrUpdateRole(serverID, newRole);
		res.match(
			() => {},
			(err) => logErr(err)
		);
	}

	setInvites(serverID: string, invites: Invite[]) {
		this.servers[serverID].invites = invites || [];
	}

	getInvites(serverID: string) {
		return this.servers[serverID].invites;
	}

	addServer(server: Server) {
		this.servers[server.id] = server;
	}

	getLastPosition() {
		return Object.values(this.servers).length;
	}

	deleteServer(serverID: string) {
		delete serverStore.servers[serverID];
		delete this.cached[serverID];
	}

	isServerInfoCached(serverID: string): boolean {
		if (this.cached[serverID]?.cached === true && this.cached[serverID].timeoutId) {
			clearTimeout(this.cached[serverID].timeoutId);
		}

		return this.cached[serverID]?.cached === true;
	}

	setCacheTimeout(serverID: string) {
		if (!this.cached[serverID]) return;
		const timeoutId = setTimeout(() => this.clearServerCache(serverID), 5 * 60 * 1000); // 5 minutes
		this.cached[serverID].timeoutId = timeoutId;
	}

	markServerInfoCached(serverID: string) {
		if (this.cached[serverID]) {
			clearTimeout(this.cached[serverID].timeoutId);
		}

		this.cached[serverID] = {
			cached: true,
			lastAccessed: Date.now()
		};
	}

	hasNotifications(serverID: string) {
		const channels = this.getServerChannels(serverID);
		const notifications = {
			unread: false,
			mentions: 0
		};

		for (const channel of channels) {
			if (channel.last_message_read !== channel.last_message_sent) {
				notifications.unread = true;
			}

			if (channel.last_mentions && channel.last_mentions?.length > 0) {
				notifications.mentions += channel.last_mentions.length;
			}
		}

		return notifications;
	}

	getServerChannels(serverID: string) {
		const server = serverStore.getServer(serverID);
		const channels = [];

		for (const category of Object.values(server.categories)) {
			channels.push(...Object.values(category.channels));
		}

		return channels;
	}

	clearServerCache(serverID: string) {
		if (this.cached[serverID]) {
			delete this.cached[serverID];
			this.servers[serverID].roles = [];
			this.servers[serverID].members = [];
			this.servers[serverID].invites = [];
			this.servers[serverID].user_roles = [];

			for (const category of Object.values(this.servers[serverID].categories)) {
				for (const channel of Object.values(category.channels)) {
					channelStore.clearChannelCache(channel.id);
				}
			}
		}
	}

	updateProfile(serverID: string, profile: EditServerType) {
		const server = this.servers[serverID];
		if (profile.name) server.name = profile.name;
		if (profile.description) server.description = profile.description;
		if (profile.public) server.public = profile.public;
	}

	updateAvatar(serverID: string, avatar?: string, banner?: string) {
		const server = this.servers[serverID];

		if (avatar) server.avatar = avatar;
		if (banner) server.banner = banner;
	}
}

export const serverStore = new ServerStore();
