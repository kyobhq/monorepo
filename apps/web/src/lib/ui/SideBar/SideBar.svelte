<script lang="ts">
	import { page } from '$app/state';
	import Input from 'ui/Input/Input.svelte';
	import MagnifyingGlass from 'ui/icons/MagnifyingGlass.svelte';
	import SideBarServers from './SideBarServers.svelte';
	import ServersSlider from 'ui/ServersSlider/ServersSlider.svelte';
	import UserBar from 'ui/UserBar/UserBar.svelte';
	import PinnedChannels from 'ui/PinnedChannels/PinnedChannels.svelte';
	import SideBarFriends from './SideBarFriends.svelte';

	const currentTab = $derived(page.url.pathname.split('/')[1]);
</script>

<aside
	class="relative bg-main-975 min-h-screen border-r-[0.5px] border-r-main-800 w-[19.5rem] overflow-hidden shrink-0"
>
	<div class="flex gap-x-2 px-2.5 pt-2.5">
		<Input Icon={MagnifyingGlass} placeholder="Search servers" />
		<!-- <IconsTabBar tabs={TABS} onclick={(href: string) => goto(href)} activeTab={page.url.pathname} /> -->
	</div>

	<ServersSlider {currentTab} />
	<PinnedChannels />

	<section class="flex flex-col relative h-[calc(100%-3.765rem)]">
		{#if currentTab === 'servers'}
			<SideBarServers />
		{:else if currentTab === 'friends'}
			<SideBarFriends />
		{/if}
	</section>
	<UserBar />
</aside>
