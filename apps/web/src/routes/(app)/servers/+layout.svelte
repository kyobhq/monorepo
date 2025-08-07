<script lang="ts">
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';

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
				serverStore.memberCount = server.member_count;
				serverStore.markServerInfoCached(page.params.server_id!);
			},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}

	let currentServerId = $state<string | undefined>();

	$effect(() => {
		const serverId = page.params.server_id;

		if (!userStore.setupComplete || !serverId) return;
		if (currentServerId === serverId || serverStore.isServerInfoCached(serverId)) return;

		currentServerId = serverId;
		const schedule = (cb: () => void) =>
			(window as any).requestIdleCallback
				? (window as any).requestIdleCallback(cb)
				: setTimeout(cb, 0);
		// Defer server info fetch to avoid blocking initial message fetch/render
		schedule(() => {
			// fire-and-forget; no need to await here
			getServerInformations();
		});
	});
</script>

{@render children()}
