<script lang="ts">
	import { page } from '$app/state';
	import { channelStore } from 'stores/channelStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import ChannelHeader from 'ui/ChannelHeader/ChannelHeader.svelte';
	import ChatBox from 'ui/ChatBox/ChatBox.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';

	const currentServer = $derived(serverStore.getServer(page.params.server_id || ''));
	const currentChannel = $derived(
		channelStore.getChannel(page.params.server_id || '', page.params.channel_id || '')
	);
</script>

{#if currentChannel && currentServer}
	<ChannelHeader name={currentChannel.name} description={currentChannel.description} />

	<div class="flex flex-col flex-1 min-h-0">
		<ChatBox serverID={currentServer.id} channelID={currentChannel.id} />

		<RichInput server={currentServer} channel={currentChannel} />
	</div>
{/if}
