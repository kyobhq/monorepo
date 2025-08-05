<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import '../app.css';
	import '@fontsource/host-grotesk/400.css';
	import '@fontsource/host-grotesk/500.css';
	import '@fontsource/host-grotesk/600.css';
	import '@fontsource/host-grotesk/700.css';
	import { backend } from 'stores/backendStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { ws } from 'stores/websocketStore.svelte';
	import { onMount } from 'svelte';
	let { children } = $props();

	onMount(async () => {
		const identityRes = await backend.getSetup();
		identityRes.match(
			(setup) => {
				userStore.user = setup.user;
				serverStore.servers = setup.servers;
				ws.init(setup.user.id);

				if (!page.params.server_id) goto('/servers');
			},
			() => goto('/signin')
		);
	});
</script>

{@render children()}
