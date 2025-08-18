<script lang="ts">
	import ContextMenuSkeleton from './ContextMenuSkeleton.svelte';
	import ContextMenuItem from './ContextMenuItem.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { generateTextWithExt } from 'utils/richInput';
	import { messageStore } from 'stores/messageStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { hasPermissions } from 'utils/permissions';

	let { message } = $props();
	const isAuthor = $derived(message.author.id === userStore.user?.id);

	async function deleteMessage() {
		const res = await backend.deleteMessage(message.id, {
			server_id: page.params.server_id || 'global',
			channel_id: page.params.channel_id!,
			author_id: message.author.id
		});

		res.match(
			() => {},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}

	async function handleDelete() {
		if (coreStore.pressingShift) {
			await deleteMessage();
			return;
		}

		coreStore.destructiveDialog = {
			open: true,
			title: `Delete Message`,
			subtitle: 'Are you sure you want to delete this message?',
			buttonText: 'Delete',
			onclick: async () => {
				coreStore.destructiveDialog.open = false;
				await deleteMessage();
			}
		};
	}

	function handleCopyText() {
		navigator.clipboard.writeText(generateTextWithExt(message.content));
	}

	function handleEdit() {
		messageStore.editMessage = message;
	}
</script>

<ContextMenuSkeleton>
	{#snippet contextMenuContent()}
		<ContextMenuItem onclick={() => {}} text="Reply" />
		<ContextMenuItem onclick={handleCopyText} text="Copy Text" />
		{#if message.server_id !== 'global'}
			<ContextMenuItem onclick={() => {}} text="Copy Message Link" />
		{/if}
		{#if isAuthor}
			<ContextMenuItem onclick={handleEdit} text="Edit Message" />
		{/if}
		{#if isAuthor || hasPermissions(page.params.server_id!, 'MANAGE_MESSAGES')}
			<ContextMenuItem onclick={handleDelete} text="Delete Message" destructive />
		{/if}
	{/snippet}
</ContextMenuSkeleton>
