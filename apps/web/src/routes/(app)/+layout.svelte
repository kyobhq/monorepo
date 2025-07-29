<script lang="ts">
	import SideBar from 'ui/SideBar/SideBar.svelte';
	import ServerBar from 'ui/ServerBar/ServerBar.svelte';
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import CreateServerDialog from 'ui/Dialog/CreateServerDialog.svelte';
	import { onMount } from 'svelte';
	import { backend } from 'stores/backendStore.svelte';
	let { children } = $props();

	const currentTab = $derived(page.url.pathname.split('/')[1]);
	const currentServer = $derived(serverStore.getServer(page.params.server_id || '') || undefined);

	onMount(async () => {
		const res = await backend.getSetup();
		res.match(
			(setup) => {
				serverStore.servers = setup.servers;
			},
			(error) => console.error(`${error.code}: ${error.message}`)
		);
	});
</script>

<div class="flex">
	<SideBar />
	<main class="flex-1 relative">
		{@render children()}
	</main>
	{#if currentTab === 'servers' && currentServer}
		<ServerBar />
	{/if}
</div>

<CreateServerDialog />
