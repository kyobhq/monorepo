<script lang="ts">
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { onMount } from 'svelte';

	let { children } = $props();

	onMount(async () => {
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
	});
</script>

{@render children()}
