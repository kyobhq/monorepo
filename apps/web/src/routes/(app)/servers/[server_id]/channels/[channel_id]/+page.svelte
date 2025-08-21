<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { onMount } from 'svelte';
	import ChannelHeader from 'ui/ChannelHeader/ChannelHeader.svelte';
	import StraightFaceEmoji from 'ui/icons/StraightFaceEmoji.svelte';
	import Message from 'ui/Message/Message.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';

	let scrollContainer = $state<HTMLDivElement>();
	let messagesLoaded = $state(false);
	let showMessages = $derived(
		coreStore.firstLoad.sidebar && coreStore.firstLoad.serverbar && messagesLoaded
	);

	const currentServer = $derived.by(() => {
		if (!page.params.server_id) return;
		return serverStore.getServer(page.params.server_id);
	});
	const currentChannel = $derived.by(() => {
		if (!page.params.server_id || !page.params.channel_id) return;
		return channelStore.getChannel(page.params.server_id, page.params.channel_id);
	});

	async function getMessages(serverID: string, channelID: string) {
		if (channelStore.messages[channelID]) {
			messagesLoaded = true;
			return;
		}

		const res = await backend.getMessages(serverID, channelID);
		res.match(
			async (messages) => {
				channelStore.messages[channelID] = messages ? messages.reverse() : [];
				messagesLoaded = true;
			},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}

	onMount(async () => {
		if (page.params.server_id && page.params.channel_id)
			await getMessages(page.params.server_id, page.params.channel_id);
	});

	beforeNavigate(async ({ from, to }) => {
		if (
			to?.params?.server_id &&
			from?.params?.channel_id &&
			to?.params?.channel_id &&
			from?.params?.channel_id !== to?.params?.channel_id
		) {
			await getMessages(to.params.server_id, to.params.channel_id);
		}
	});
</script>

{#if currentChannel && currentServer}
	<ChannelHeader name={currentChannel.name} description={currentChannel.description} />

	<div
		class="flex flex-col-reverse w-full h-full overflow-auto pb-4 pt-18"
		bind:this={scrollContainer}
	>
		{#if showMessages}
			{#if channelStore.messages[currentChannel.id].length > 0}
				{#each channelStore.messages[currentChannel.id] as message (message.id)}
					<Message server={currentServer} channel={currentChannel} {message} />
				{/each}
			{:else}
				<div
					class="h-full w-full flex items-center justify-center text-2xl flex-col gap-y-4 text-main-700"
				>
					<StraightFaceEmoji height={128} width={128} />
					<p>No messages in <span class="font-medium">{currentChannel?.name}</span> yet</p>
				</div>
			{/if}
		{/if}
	</div>

	<RichInput server={currentServer} channel={currentChannel} />
{/if}
