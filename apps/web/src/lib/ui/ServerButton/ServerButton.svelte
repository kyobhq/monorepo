<script lang="ts">
	import type { Server } from '$lib/types/types';
	import ContextMenuServer from 'ui/ContextMenu/ContextMenuServer.svelte';

	interface Props {
		server: Server;
		active?: boolean;
		onclick: () => void;
	}

	let { server, active, onclick }: Props = $props();
</script>

<button
	class={[
		'h-13 w-13 aspect-square object-cover relative after:content-normal after:absolute after:inset-0 after:transition after:pointer-events-none after:rounded-[inherit] transition-all active-scale-down',
		active ? 'after:inner-active rounded-lg' : 'after:inner-main-700 rounded-xl hover:rounded-lg'
	]}
	{onclick}
>
	<img src={server.avatar} alt="server-icon" class="h-full w-full rounded-[inherit]" />
	<ContextMenuServer {server} />
</button>
