<script lang="ts">
	import { page } from '$app/state';
	import type { Server } from '$lib/types/types';
	import { serverStore } from 'stores/serverStore.svelte';
	import Avatar from 'ui/Avatar/Avatar.svelte';
	import ContextMenuServer from 'ui/ContextMenu/ContextMenuServer.svelte';

	interface Props {
		server: Server;
		active?: boolean;
		onclick: () => void;
	}

	let { server, active, onclick }: Props = $props();
	let hoverAvatar = $state(false);

	const hasNotifications = $derived(
		page.params.server_id !== server.id && serverStore.hasNotifications(server.id)
	);
</script>

<button
	class={[
		'h-13 w-13 aspect-square object-cover relative after:content-normal after:absolute after:inset-0 after:transition after:pointer-events-none after:rounded-[inherit] transition-all active-scale-down',
		active ? 'after:inner-active rounded-lg' : 'after:inner-main-700 rounded-xl hover:rounded-lg'
	]}
	{onclick}
	onmouseenter={() => (hoverAvatar = true)}
	onmouseleave={() => (hoverAvatar = false)}
>
	<Avatar
		src={server.avatar}
		alt="server-icon"
		class="h-full w-full rounded-[inherit] overflow-hidden"
		hover={hoverAvatar}
	/>
	<ContextMenuServer {server} />

	{#if hasNotifications}
		{#if hasNotifications.mentions}
			<div
				class="absolute -top-1 -right-1 w-5 h-5 bg-red-400 rounded-lg flex items-center justify-center text-xs font-semibold z-[1]"
			>
				{hasNotifications.mentions}
			</div>
		{/if}
		{#if hasNotifications.unread}
			<div
				class="absolute -bottom-0.5 bg-main-50 left-1/2 -translate-x-1/2 w-4 h-1 rounded-full"
			></div>
		{/if}
	{/if}
</button>
