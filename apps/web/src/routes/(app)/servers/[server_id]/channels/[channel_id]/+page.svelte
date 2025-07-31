<script lang="ts">
	import { page } from '$app/state';
	import { channelStore } from 'stores/channelStore.svelte';
	import ChannelHeader from 'ui/ChannelHeader/ChannelHeader.svelte';

	const currentChannel = $derived.by(() => {
		if (!page.params.server_id || !page.params.channel_id) return;
		return channelStore.getChannel(page.params.server_id, page.params.channel_id);
	});
</script>

{#if currentChannel}
	<ChannelHeader name={currentChannel.name} description={currentChannel.description} />
{/if}
