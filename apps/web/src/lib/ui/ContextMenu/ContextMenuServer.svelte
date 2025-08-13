<script lang="ts">
	import ContextMenuSkeleton from './ContextMenuSkeleton.svelte';
	import ContextMenuItem from './ContextMenuItem.svelte';
	import type { Server } from '$lib/types/types';
	import { userStore } from 'stores/userStore.svelte';
	import { backend } from 'stores/backendStore.svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { serverStore } from 'stores/serverStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { hasPermissions } from 'utils/permissions';

	interface Props {
		server: Server;
	}

	let { server }: Props = $props();
	let inviteCreated = $state(false);

	function openSettings() {
		coreStore.serverSettingsDialog = {
			open: true,
			server_id: server.id,
			section: 'Server Profile'
		};
	}

	async function createInvite(e: Event) {
		e.preventDefault();

		const res = await backend.getInviteLink(server.id);
		res.match(
			(inviteLink) => {
				navigator.clipboard.writeText(inviteLink);
				inviteCreated = true;
				setTimeout(() => {
					inviteCreated = false;
				}, 1500);
			},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}

	async function leaveServer() {
		const res = await backend.leaveServer(server.id);
		res.match(
			() => {
				serverStore.deleteServer(server.id);

				if (page.params.server_id === server.id) {
					goto('/servers');
				}
			},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}
</script>

<ContextMenuSkeleton>
	{#snippet contextMenuContent()}
		{#if hasPermissions(server.id, 'CREATE_INVITE')}
			<ContextMenuItem
				onclick={createInvite}
				text={inviteCreated ? 'Check your clipboard!' : 'Get Invite Link'}
				success={inviteCreated}
			/>
		{/if}
		{#if hasPermissions(server.id, 'MANAGE_SERVER', 'MANAGE_ROLES')}
			<ContextMenuItem onclick={openSettings} text="Edit Server" />
		{/if}
		{#if server.owner_id === userStore.user!.id}
			<ContextMenuItem onclick={() => {}} text="Delete Server" destructive />
		{:else}
			<ContextMenuItem onclick={leaveServer} text="Leave Server" destructive />
		{/if}
	{/snippet}
</ContextMenuSkeleton>
