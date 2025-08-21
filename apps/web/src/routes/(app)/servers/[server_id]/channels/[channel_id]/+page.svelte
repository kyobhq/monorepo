<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import { channelStore } from 'stores/channelStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { onMount } from 'svelte';
	import ChannelHeader from 'ui/ChannelHeader/ChannelHeader.svelte';
	import StraightFaceEmoji from 'ui/icons/StraightFaceEmoji.svelte';
	import Message from 'ui/Message/Message.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';
	import { delay } from 'utils/time';

	const SCROLL_THRESHOLD = 500;

	const currentServer = $derived(serverStore.getServer(page.params.server_id || ''));
	const currentChannel = $derived(
		channelStore.getChannel(page.params.server_id || '', page.params.channel_id || '')
	);

	let messagesLoaded = $state(false);
	let canLoadMore = $state(true);
	let showMessages = $derived(
		coreStore.firstLoad.sidebar && coreStore.firstLoad.serverbar && messagesLoaded
	);

	async function handleScroll(e: Event) {
		if (!currentChannel || !currentServer || !canLoadMore) return;
		const target = e.target as HTMLDivElement;
		const scrollY = target.clientHeight + Math.abs(target.scrollTop);

		if (target.scrollHeight - SCROLL_THRESHOLD <= scrollY) {
			canLoadMore = false;

			const hasMore = await channelStore.loadMoreMessages(
				currentServer.id,
				currentChannel.id,
				'before'
			);

			if (hasMore) {
				setTimeout(() => {
					canLoadMore = true;
				}, 500);
			}
		} else if (
			channelStore.messageCache[currentChannel.id]?.offsetHeight > 0 &&
			scrollY - channelStore.messageCache[currentChannel.id]?.offsetHeight - target.clientHeight <=
				SCROLL_THRESHOLD
		) {
			canLoadMore = false;

			const hasMore = await channelStore.loadMoreMessages(
				currentServer.id,
				currentChannel.id,
				'after'
			);

			if (hasMore) {
				setTimeout(() => {
					canLoadMore = true;
				}, 500);
			}
		}
	}

	onMount(async () => {
		if (page.params.server_id && page.params.channel_id) {
			await channelStore.ensureMessagesLoaded(page.params.server_id, page.params.channel_id);
			messagesLoaded = true;
		}
	});

	beforeNavigate(async ({ from, to }) => {
		messagesLoaded = false;
		const fromChannelID = from?.params?.channel_id;
		const toChannelID = to?.params?.channel_id;
		if (!fromChannelID || !toChannelID) return;

		const serverID = from?.params?.server_id || to?.params?.server_id;
		const sameChannel = fromChannelID === toChannelID;
		if (serverID && !sameChannel) {
			await channelStore.ensureMessagesLoaded(serverID, toChannelID);
			canLoadMore = !(channelStore.messageCache[toChannelID]?.hasReachedEnd ?? false);
		}

		await delay(0);
		messagesLoaded = true;
	});
</script>

{#if currentChannel && currentServer}
	<ChannelHeader name={currentChannel.name} description={currentChannel.description} />

	<div class="flex flex-col-reverse w-full h-full overflow-auto pb-4 pt-18" onscroll={handleScroll}>
		{#if showMessages}
			{#if channelStore.messageCache[currentChannel.id]?.messages.length > 0}
				{#if channelStore.messageCache[currentChannel.id]?.offsetHeight > 0}
					<div
						style="height: {channelStore.messageCache[currentChannel.id]
							.offsetHeight}px; flex-shrink: 0;"
						aria-hidden="true"
					></div>
				{/if}

				{#each channelStore.messageCache[currentChannel.id].messages as message (message.id)}
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
