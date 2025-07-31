<script lang="ts">
	import { page } from '$app/state';
	import { channelStore } from 'stores/channelStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import ChannelHeader from 'ui/ChannelHeader/ChannelHeader.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';

	const currentServer = $derived.by(() => {
		if (!page.params.server_id) return;
		return serverStore.getServer(page.params.server_id);
	});
	const currentChannel = $derived.by(() => {
		if (!page.params.server_id || !page.params.channel_id) return;
		return channelStore.getChannel(page.params.server_id, page.params.channel_id);
	});
</script>

{#if currentChannel}
	<ChannelHeader name={currentChannel.name} description={currentChannel.description} />
{/if}

<div class="flex w-full h-full"></div>

{#if currentServer && currentChannel}
	<RichInput server={currentServer} channel={currentChannel} />
{/if}
