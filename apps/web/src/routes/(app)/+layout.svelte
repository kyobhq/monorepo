<script lang="ts">
	import SideBar from 'ui/SideBar/SideBar.svelte';
	import ServerBar from 'ui/ServerBar/ServerBar.svelte';
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import ServerDialog from 'ui/Dialog/ServerDialog/ServerDialog.svelte';
	import CreateCategoryDialog from 'ui/Dialog/CreateCategoryDialog/CreateCategoryDialog.svelte';
	import CreateChannelDialog from 'ui/Dialog/CreateChannelDialog/CreateChannelDialog.svelte';
	import DestructiveDialog from 'ui/Dialog/DestructiveDialog/DestructiveDialog.svelte';
	import ChannelSettingsDialog from 'ui/Dialog/ChannelSettingsDialog/ChannelSettingsDialog.svelte';
	import { onDestroy, onMount } from 'svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import UserProfile from 'ui/UserProfile/UserProfile.svelte';
	import UserSettingsDialog from 'ui/Dialog/UserSettingsDialog/UserSettingsDialog.svelte';
	import ServerSettingsDialog from 'ui/Dialog/ServerSettingsDialog/ServerSettingsDialog.svelte';
	import FriendsDialog from 'ui/Dialog/FriendsDialog/FriendsDialog.svelte';
	import { onNavigate } from '$app/navigation';
	import CategorySettingsDialog from 'ui/Dialog/CategorySettingsDialog/CategorySettingsDialog.svelte';

	let { children } = $props();

	const currentTab = $derived(page.url.pathname.split('/')[1]);
	const currentServer = $derived(serverStore.getServer(page.params.server_id || '') || undefined);
	let mainEl = $state<HTMLElement>();

	function handleRichInputLength() {
		if (mainEl) coreStore.richInputLength = mainEl.clientWidth;
	}

	onMount(() => {
		coreStore.initializeKeyboardDetection();
		handleRichInputLength();
		window.addEventListener('resize', handleRichInputLength);
	});

	onDestroy(() => {
		coreStore.cleanupKeyboardDetection();
		window.removeEventListener('resize', handleRichInputLength);
	});

	onNavigate(() => {
		handleRichInputLength();
	});
</script>

<div class="flex">
	<SideBar />
	<main
		bind:this={mainEl}
		class={[
			'flex flex-col h-screen relative',
			page.url.pathname.includes('servers') ? 'w-[calc(100%-19.5rem-16rem)]' : 'w-full'
		]}
	>
		{@render children()}
	</main>
	{#if currentTab === 'servers' && currentServer}
		<ServerBar />
	{/if}
</div>

<FriendsDialog />
<ServerDialog />
<CreateCategoryDialog />
<CreateChannelDialog />
<DestructiveDialog />
<ChannelSettingsDialog />
<CategorySettingsDialog />
<UserSettingsDialog />
<ServerSettingsDialog />
<UserProfile />
