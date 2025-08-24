import type { Category } from '$lib/types/types';
import { serverStore } from './serverStore.svelte';

class CategoryStore {
	getCategory(serverID: string, categoryID: string): Category | null {
		return serverStore.safeServerOperation(
			serverID,
			(server) => server.categories[categoryID] || null,
			null
		);
	}

	getLastPositionInServer(serverID: string): number {
		return serverStore.safeServerOperation(
			serverID,
			(server) => Object.values(server.categories).length,
			0
		);
	}

	addCategory(category: Category): void {
		serverStore.safeServerOperation(
			category.server_id,
			(server) => {
				server.categories[category.id] = category;
				return true;
			},
			false
		);
	}

	editCategory(categoryID: string, categoryOpts: Partial<Category>): void {
		const category = this.getCategory(categoryOpts.server_id!, categoryID);
		if (!category) return;

		if (categoryOpts.name) category.name = categoryOpts.name;
		if (categoryOpts.users) category.users = categoryOpts.users;
		if (categoryOpts.roles) category.roles = categoryOpts.roles;
	}

	deleteCategory(serverID: string, categoryID: string): void {
		serverStore.safeServerOperation(
			serverID,
			(server) => {
				delete server.categories[categoryID];
				return true;
			},
			false
		);
	}
}

export const categoryStore = new CategoryStore();
