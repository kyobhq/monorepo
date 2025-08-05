<script lang="ts">
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';

	let { children } = $props();

	async function getServerInformations() {
		if (!page.params.server_id) return;

		const serverInformations = await backend.getServerInformations(page.params.server_id);
		serverInformations.match(
			(server) => {
				serverStore.members = server.members;
				serverStore.roles = server.roles;
				serverStore.memberCount = server.member_count;
			},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}

	$effect(() => {
		if (!page.params.server_id) return;
		getServerInformations();
	});
</script>

{@render children()}
