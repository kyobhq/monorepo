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

  addCategory(category: Category) {
    serverStore.safeServerOperation(category.server_id, (server) => {
      server.categories[category.id] = category
      return true
    }, false)
  }

  editCategory(categoryID: string, categoryOpts: Partial<Category>) {
    const category = this.getCategory(categoryOpts.server_id!, categoryID)
    if (!category) return;

    if (categoryOpts.name) category.name = categoryOpts.name
    if (categoryOpts.users) category.users = categoryOpts.users
    if (categoryOpts.roles) category.roles = categoryOpts.roles
  }

  deleteCategory(serverID: string, categoryID: string) {
    serverStore.safeServerOperation(serverID, (server) => {
      delete server.categories[categoryID]
      return true
    }, false)
  }
}

export const categoryStore = new CategoryStore();
