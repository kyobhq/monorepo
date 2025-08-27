<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { channelStore } from 'stores/channelStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import ChannelHeader from 'ui/ChannelHeader/ChannelHeader.svelte';
	import ChatBox from 'ui/ChatBox/ChatBox.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';

	const currentServer = $derived(serverStore.getServer(page.params.server_id || ''));
	const currentChannel = $derived(
		channelStore.getChannel(page.params.server_id || '', page.params.channel_id || '')
	);

	$effect(() => {
		if (
			currentChannel?.users &&
			currentChannel.users.length > 0 &&
			!currentChannel.users.includes(userStore.user!.id)
		) {
			goto(`/servers/${currentServer.id}`);
		}
	});
</script>

{#if currentChannel && currentServer}
	<ChannelHeader name={currentChannel.name} description={currentChannel.description} />

	<div class="flex flex-col flex-1 min-h-0">
		<ChatBox serverID={currentServer.id} channelID={currentChannel.id} />

		<RichInput server={currentServer} channel={currentChannel} />
	</div>
{/if}
