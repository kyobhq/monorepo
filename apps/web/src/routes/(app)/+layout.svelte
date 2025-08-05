<script lang="ts">
	import SideBar from 'ui/SideBar/SideBar.svelte';
	import ServerBar from 'ui/ServerBar/ServerBar.svelte';
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import CreateServerDialog from 'ui/Dialog/CreateServerDialog.svelte';
	import CreateCategoryDialog from 'ui/Dialog/CreateCategoryDialog.svelte';
	import CreateChannelDialog from 'ui/Dialog/CreateChannelDialog.svelte';
	import DestructiveDialog from 'ui/Dialog/DestructiveDialog.svelte';
	import ChannelSettingsDialog from 'ui/Dialog/ChannelSettingsDialog.svelte';
	import { onDestroy, onMount } from 'svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import UserProfile from 'ui/UserProfile/UserProfile.svelte';
	let { children } = $props();

	const currentTab = $derived(page.url.pathname.split('/')[1]);
	const currentServer = $derived(serverStore.getServer(page.params.server_id || '') || undefined);

	onMount(() => {
		window.addEventListener('keydown', coreStore.handleShiftDown);
		window.addEventListener('keyup', coreStore.handleShiftUp);
	});

	onDestroy(() => {
		window.removeEventListener('keydown', coreStore.handleShiftDown);
		window.removeEventListener('keyup', coreStore.handleShiftUp);
	});
</script>

<div class="flex">
	<SideBar />
	<main class="flex flex-col w-[calc(100%-19.5rem*2)] h-screen relative">
		{@render children()}
	</main>
	{#if currentTab === 'servers' && currentServer}
		<ServerBar />
	{/if}
</div>

<CreateServerDialog />
<CreateCategoryDialog />
<CreateChannelDialog />
<DestructiveDialog />
<ChannelSettingsDialog />
<UserProfile />
