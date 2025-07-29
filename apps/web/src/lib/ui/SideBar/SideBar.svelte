<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import IconsTabBar from 'ui/IconsTabBar/IconsTabBar.svelte';
	import Input from 'ui/Input/Input.svelte';
	import MagnifyingGlass from 'ui/icons/MagnifyingGlass.svelte';
	import { TABS } from '../../constants/constants';
	import SideBarServers from './SideBarServers.svelte';
	import SideBarFriends from './SideBarFriends.svelte';

	const currentTab = $derived(page.url.pathname.split('/')[1]);
</script>

<aside
	class="bg-main-975 min-h-screen border-r-[0.5px] border-r-main-700 w-[19.5rem] overflow-hidden"
>
	<div class="flex gap-x-2 px-2.5 pt-2.5">
		<Input Icon={MagnifyingGlass} placeholder="Search" />
		<IconsTabBar tabs={TABS} onclick={(href: string) => goto(href)} activeTab={page.url.pathname} />
	</div>
	{#if currentTab === 'servers'}
		<SideBarServers />
	{:else}
		<SideBarFriends />
	{/if}
</aside>
