import type { Server } from "$lib/types/types";
import { print } from "utils/print"

export class ServerStore {
  servers = $state<Record<string, Server>>({});

  getChannel() {
    print("todo!")
  }

  getServer(id: string) {
    return this.servers[id]
  }

  addServer(server: Server) {
    this.servers[server.id] = server
  }
}

export const serverStore = new ServerStore()
