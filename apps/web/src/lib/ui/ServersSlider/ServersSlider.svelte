<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import People from 'ui/icons/People.svelte';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';
	import ServerButton from 'ui/ServerButton/ServerButton.svelte';

	let { currentTab } = $props();

	const defaultButton =
		'h-13 w-13 bg-main-900 text-main-500 aspect-square hover:bg-main-800 hover:text-main-200 hover:cursor-pointer transition-colors duration-75 flex items-center justify-center rounded-xl active-scale-down';
</script>

<div
	class="flex gap-x-2.5 py-2.5 pl-2.5 shrink-0 relative after:absolute after:top-0 after:right-0 after:w-16 after:h-full after:bg-gradient-to-l after:from-main-975 after:to-transparent after:pointer-events-none"
>
	<button
		class={[defaultButton, currentTab === 'friends' && 'bg-accent! text-main-50!']}
		onclick={() => goto(`/friends`)}
	>
		<People height={20} width={20} />
	</button>
	{#each Object.values(serverStore.servers).sort((a, b) => a.position - b.position) as server (server.id)}
		<ServerButton
			{server}
			onclick={() => goto(`/servers/${server.id}`)}
			active={page.url.pathname.includes(server.id)}
		/>
	{/each}
	<button class={defaultButton} onclick={() => (coreStore.serverDialog = true)}>
		<PlusIcon height={20} width={20} />
	</button>
	{#if Object.keys(serverStore.servers).length < 4}
		{#each { length: 4 - Object.keys(serverStore.servers).length }, _}
			<div class="h-13 w-13 bg-main-950 aspect-square rounded-xl"></div>
		{/each}
	{/if}
</div>
