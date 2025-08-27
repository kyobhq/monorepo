<script lang="ts">
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import ChannelsList from 'ui/ChannelsList/ChannelsList.svelte';
	import BarSeparator from 'ui/BarSeparator/BarSeparator.svelte';

	const currentServer = $derived(serverStore.getServer(page.params.server_id || '') || undefined);
</script>

{#if currentServer}
	<BarSeparator title={currentServer.name} />
	<ChannelsList server={currentServer} />
{/if}
