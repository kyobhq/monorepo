import type { Member, Role, Server } from '$lib/types/types';

export class ServerStore {
	servers = $state<Record<string, Server>>({});
	memberCount = $state<number>(0);
	cached = $state<Record<string, boolean>>({});

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

	getMember(serverID: string, userID: string) {
		return this.servers[serverID].members.find((member) => member.id === userID);
	}

	setMembers(serverID: string, members: Member[]) {
		this.servers[serverID].members = members;
	}

	getMembers(serverID: string) {
		return this.servers[serverID].members;
	}

	setRoles(serverID: string, roles: Role[]) {
		this.servers[serverID].roles = roles || [];
	}

	getRoles(serverID: string) {
		return this.servers[serverID].roles;
	}

	addServer(server: Server) {
		this.servers[server.id] = server;
	}

	getLastPosition() {
		return Object.values(this.servers).length;
	}

	addMember(serverID: string, member: Member) {
		this.servers[serverID].members.push(member);
	}

	setMemberOnline(serverID: string, memberID: string, status: string) {
		const member = this.servers[serverID].members.find((member) => member.id === memberID);
		if (member) member.status = status;
	}

	setMemberOffline(serverID: string, memberID: string) {
		const member = this.servers[serverID].members.find((member) => member.id === memberID);
		if (member) member.status = 'offline';
	}

	deleteServer(serverID: string) {
		delete serverStore.servers[serverID];
		delete this.cached[serverID];
	}

	isServerInfoCached(serverID: string): boolean {
		return this.cached[serverID] === true;
	}

	markServerInfoCached(serverID: string) {
		this.cached[serverID] = true;
	}
}

export const serverStore = new ServerStore();
