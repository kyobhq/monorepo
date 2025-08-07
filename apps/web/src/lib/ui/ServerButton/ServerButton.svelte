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
		'h-13 w-13 aspect-square object-cover relative after:content-normal after:absolute after:inset-0 after:transition after:pointer-events-none',
		active
			? 'after:inner-active'
			: 'after:inner-main-700 hocus:after:inner-active hover:cursor-pointer'
	]}
	{onclick}
>
	<img src={server.avatar} alt="server-icon" class="h-full w-full" />
	<ContextMenuServer {server} />
</button>
