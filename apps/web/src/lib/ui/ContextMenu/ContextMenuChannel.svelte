<script lang="ts">
	import ContextMenuSkeleton from './ContextMenuSkeleton.svelte';
	import ContextMenuItem from './ContextMenuItem.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { hasPermissions } from 'utils/permissions';

	let { categoryId, channelId } = $props();

	function handleDelete() {
		if (!page.params.server_id) return;

		const channel = channelStore.getChannel(page.params.server_id, channelId);
		if (!channel) return;

		coreStore.destructiveDialog = {
			open: true,
			title: `Delete ${channel.name}`,
			subtitle: 'All content in this channel will be permanently deleted.',
			buttonText: 'Delete channel',
			onclick: async () => {
				coreStore.destructiveDialog.open = false;
				const res = await backend.deleteChannel(page.params.server_id!, categoryId, channelId);
				res.match(
					() => {},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);
			}
		};
	}

	function openSettings() {
		coreStore.channelSettingsDialog = {
			open: true,
			channel_id: channelId,
			section: 'Overview'
		};
	}
</script>

<ContextMenuSkeleton>
	{#snippet contextMenuContent()}
		<ContextMenuItem onclick={() => {}} text="Copy Link" />
		{#if hasPermissions(page.params.server_id!, 'MANAGE_CHANNELS')}
			<ContextMenuItem onclick={openSettings} text="Edit Channel" />
			<ContextMenuItem onclick={handleDelete} text="Delete Channel" destructive />
		{/if}
	{/snippet}
</ContextMenuSkeleton>
