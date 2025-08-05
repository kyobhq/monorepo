import type { Member, Role, Server } from "$lib/types/types";

export class ServerStore {
  servers = $state<Record<string, Server>>({});
  members = $state<Member[]>([]);
  roles = $state<Role[]>([]);
  memberCount = $state<number>(0);

  safeServerOperation<T>(
    serverID: string,
    operation: (server: Server) => T,
    fallback: T
  ): T {
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
    return this.servers[id]
  }

  addServer(server: Server) {
    this.servers[server.id] = server
  }

  getLastPosition() {
    return Object.values(this.servers).length
  }

  deleteServer(serverID: string) {
    delete serverStore.servers[serverID]
  }
}

export const serverStore = new ServerStore()
