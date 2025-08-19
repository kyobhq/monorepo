<script lang="ts">
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import TopGradient from 'ui/TopGradient/TopGradient.svelte';

	let { children } = $props();

	async function getServerInformations() {
		if (!page.params.server_id) return;

		const serverExists = serverStore.getServer(page.params.server_id);
		if (!serverExists) {
			console.warn(`Server ${page.params.server_id} not found in store`);
			return;
		}

		const serverInformations = await backend.getServerInformations(page.params.server_id);
		serverInformations.match(
			(server) => {
				serverStore.setMembers(page.params.server_id!, server.members);
				serverStore.setRoles(page.params.server_id!, server.roles);
				serverStore.setUserRoles(page.params.server_id!, server.user_roles);
				serverStore.setInvites(page.params.server_id!, server.invites);
				serverStore.memberCount = server.member_count;
				serverStore.markServerInfoCached(page.params.server_id!);
			},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}

	let currentServerId = $state<string | undefined>();
	let currentServerName = $derived(
		page.params?.server_id && serverStore.getServer(page.params.server_id)?.name
	);

	$effect(() => {
		const serverId = page.params.server_id;

		if (!userStore.setupComplete || !serverId) return;
		if (currentServerId === serverId || serverStore.isServerInfoCached(serverId)) return;

		currentServerId = serverId;
		getServerInformations();
	});
</script>

<svelte:head>
	<title>Kyob | {currentServerName}</title>
</svelte:head>

{@render children()}

<TopGradient />
