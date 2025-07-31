import type { Category } from "$lib/types/types";
import { serverStore } from "./serverStore.svelte";

class CategoryStore {
  getCategory(serverID: string, categoryID: string) {
    return serverStore.safeServerOperation(serverID, (server) => {
      return server.categories[categoryID] || null
    }, null)
  }

  getLastPositionInServer(serverID: string) {
    return serverStore.safeServerOperation(serverID, (server) => {
      return Object.values(server.categories).length || 0
    }, 0)
  }

  addCategory(serverID: string, category: Category) {
    serverStore.safeServerOperation(serverID, (server) => {
      server.categories[category.id] = category
      return true
    }, false)
  }

  deleteCategory(serverID: string, categoryID: string) {
    serverStore.safeServerOperation(serverID, (server) => {
      delete server.categories[categoryID]
      return true
    }, false)
  }
}

export const categoryStore = new CategoryStore();
