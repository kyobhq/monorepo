<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';
	import ServerButton from 'ui/ServerButton/ServerButton.svelte';
</script>

<div
	class="flex gap-x-2.5 pl-2.5 pb-2.5 shrink-0 relative after:absolute after:top-0 after:right-0 after:w-16 after:h-full after:bg-gradient-to-l after:from-main-975 after:to-transparent after:pointer-events-none"
>
	<button
		class="h-13 w-13 bg-main-900 text-main-500 aspect-square hocus:bg-main-800 hocus:text-main-200 hover:cursor-pointer transition-colors duration-75 flex items-center justify-center"
		onclick={() => (coreStore.openServerDialog = true)}
	>
		<PlusIcon height={20} width={20} />
	</button>
	{#each Object.values(serverStore.servers).sort((a, b) => a.position - b.position) as server (server.id)}
		<ServerButton
			image={server.avatar}
			onclick={() => goto(`/servers/${server.id}`)}
			active={page.url.pathname.includes(server.id)}
		/>
	{/each}
	{#each { length: 5 }, _}
		<div class="h-13 w-13 bg-main-950 aspect-square"></div>
	{/each}
</div>
