<script lang="ts">
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import BarSeparatorServers from 'ui/BarSeparator/BarSeparatorServers.svelte';
	import MembersList from 'ui/MembersList/MembersList.svelte';
	import ChannelsList from 'ui/ChannelsList/ChannelsList.svelte';

	const currentServer = $derived(serverStore.getServer(page.params.server_id || '') || undefined);

	let serverTab = $state<'channels' | 'members'>('channels');
</script>

{#if currentServer}
	<BarSeparatorServers title={currentServer.name} bind:tab={serverTab} />
	{#if serverTab === 'channels'}
		<ChannelsList server={currentServer} />
	{:else if serverTab === 'members'}
		<MembersList />
	{/if}
{/if}
