<script lang="ts">
	import ContextMenuSkeleton from './ContextMenuSkeleton.svelte';
	import ContextMenuItem from './ContextMenuItem.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';

	let { channelId } = $props();

	function handleDelete() {
		if (!page.params.server_id) return;

		const channel = channelStore.getChannel(page.params.server_id, channelId);
		if (!channel) return;

		coreStore.openDestructiveDialog = {
			open: true,
			title: `Delete ${channel.name}`,
			subtitle: 'All content in this channel will be permanently deleted.',
			buttonText: 'Delete channel',
			onclick: async () => {
				const res = await backend.deleteChannel(page.params.server_id!, channelId);
				res.match(
					(_) => {
						coreStore.openDestructiveDialog.open = false;
					},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);
			}
		};
	}
</script>

<ContextMenuSkeleton>
	{#snippet contextMenuContent()}
		<ContextMenuItem onclick={() => {}} text="Edit Channel" />
		<ContextMenuItem onclick={handleDelete} text="Delete Channel" destructive />
	{/snippet}
</ContextMenuSkeleton>
