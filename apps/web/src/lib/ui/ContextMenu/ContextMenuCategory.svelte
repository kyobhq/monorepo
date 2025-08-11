<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import ContextMenuSkeleton from './ContextMenuSkeleton.svelte';
	import ContextMenuItem from './ContextMenuItem.svelte';
	import { categoryStore } from 'stores/categoryStore.svelte';
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';

	let { categoryId } = $props();
	let currentServer = $derived(serverStore.getServer(page.params.server_id!));

	function openChannelDialog() {
		coreStore.channelDialog = {
			open: true,
			category_id: categoryId
		};
	}

	function handleDelete() {
		if (!page.params.server_id) return;

		const category = categoryStore.getCategory(page.params.server_id, categoryId);
		if (!category) return;

		coreStore.destructiveDialog = {
			open: true,
			title: `Delete ${category.name}`,
			subtitle: 'All content in the channels will be permanently deleted.',
			buttonText: 'Delete category',
			onclick: async () => {
				coreStore.destructiveDialog.open = false;
				const res = await backend.deleteCategory(page.params.server_id!, categoryId);
				res.match(
					() => {},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);
			}
		};
	}
</script>

{#if currentServer.owner_id === userStore.user?.id}
	<ContextMenuSkeleton>
		{#snippet contextMenuContent()}
			<ContextMenuItem onclick={openChannelDialog} text="Create channel" />
			<ContextMenuItem onclick={() => {}} text="Edit category" />
			<ContextMenuItem onclick={handleDelete} text="Delete category" destructive />
		{/snippet}
	</ContextMenuSkeleton>
{/if}
