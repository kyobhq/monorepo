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
	import { logErr } from 'utils/print';

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
		coreStore.destructiveDialog = {
			open: true,
			title: `Leave ${server.name}`,
			subtitle: server.public
				? "You can always join it again from the discovery tab don't worry."
				: "You'll need a valid invite to join this server again.",
			buttonText: 'Leave Server',
			onclick: async () => {
				coreStore.destructiveDialog.open = false;
				const res = await backend.leaveServer(server.id);
				res.match(
					() => {
						if (page.params.server_id === server.id) goto('/servers');
						serverStore.deleteServer(server.id);
					},
					(error) => logErr(error)
				);
			}
		};
	}

	async function deleteServer() {
		coreStore.destructiveDialog = {
			open: true,
			title: `Delete ${server.name}`,
			subtitle: 'All content in this server will be permanently deleted.',
			buttonText: 'Delete Server',
			onclick: async () => {
				coreStore.destructiveDialog.open = false;
				const res = await backend.deleteServer(server.id);
				if (res.isErr()) logErr(res.error);
			}
		};
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
			<ContextMenuItem onclick={deleteServer} text="Delete Server" destructive />
		{:else}
			<ContextMenuItem onclick={leaveServer} text="Leave Server" destructive />
		{/if}
	{/snippet}
</ContextMenuSkeleton>
