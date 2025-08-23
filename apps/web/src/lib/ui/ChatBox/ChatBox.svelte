<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import { channelStore } from 'stores/channelStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { onMount } from 'svelte';
	import Message from 'ui/Message/Message.svelte';
	import { delay } from 'utils/time';
	import ChatSpacer from './ChatSpacer.svelte';
	import ChatNoMessages from './ChatNoMessages.svelte';
	import { fade } from 'svelte/transition';

	interface Props {
		serverID: string;
		channelID: string;
	}

	let { serverID, channelID }: Props = $props();

	const currentServer = $derived(serverStore.getServer(serverID));
	const currentChannel = $derived(channelStore.getChannel(serverID, channelID));

	const SCROLL_THRESHOLD = 500;
	const SCROLL_LIMIT = 3000;
	let scrollSet = $state(false);
	let scrollY = $state(0);
	let scrollContainer = $state<HTMLDivElement>();
	let messagesLoaded = $state(false);
	let canLoadMore = $state(true);
	let showMessages = $derived(coreStore.firstLoad.sidebar && messagesLoaded);

	async function handleScroll(e: Event) {
		if (!canLoadMore || !scrollSet) return;
		const target = e.target as HTMLDivElement;
		scrollY = target.clientHeight + Math.abs(target.scrollTop);

		if (target.scrollHeight - SCROLL_THRESHOLD <= scrollY) {
			canLoadMore = false;

			const hasMore = await channelStore.loadMoreMessages(serverID, channelID, 'before');

			if (hasMore) {
				setTimeout(() => {
					canLoadMore = true;
				}, 500);
			}
		} else if (
			channelStore.messageCache[channelID]?.scrollHeight > 0 &&
			scrollY - channelStore.messageCache[channelID]?.scrollHeight - target.clientHeight <=
				SCROLL_THRESHOLD
		) {
			canLoadMore = false;

			const hasMore = await channelStore.loadMoreMessages(serverID, channelID, 'after');

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
			await delay(100);
			messagesLoaded = true;
		}
	});

	beforeNavigate(async ({ from, to }) => {
		scrollSet = false;
		messagesLoaded = false;
		const fromChannelID = from?.params?.channel_id;
		const toChannelID = to?.params?.channel_id;
		if (!fromChannelID || !toChannelID) return;

		if (scrollContainer) {
			channelStore.messageCache[fromChannelID].scrollY = scrollContainer.scrollTop;
			scrollContainer.scrollTop = 0;
		}

		const serverID = from?.params?.server_id || to?.params?.server_id;
		const sameChannel = fromChannelID === toChannelID;
		if (serverID && !sameChannel) {
			await channelStore.ensureMessagesLoaded(serverID, toChannelID);
			canLoadMore = !(channelStore.messageCache[toChannelID]?.hasReachedEnd ?? false);
		}

		await delay(0);
		messagesLoaded = true;
	});

	$effect(() => {
		if (scrollContainer && messagesLoaded) {
			scrollContainer.scrollTop = channelStore.messageCache[channelID]?.scrollY ?? 0;
			scrollSet = true;
		}
	});

	$effect(() => {
		if (scrollContainer && channelStore.messageCache[channelID]?.messages.length) {
			if (Math.abs(scrollContainer.scrollTop) < SCROLL_LIMIT) scrollContainer.scrollTop = 0;
		}
	});
</script>

<div class="w-full h-[calc(100%-4.5rem)]">
	{#if showMessages}
		<div
			transition:fade={{ duration: 100 }}
			bind:this={scrollContainer}
			class="flex flex-col-reverse w-full h-full overflow-auto pb-4 pt-18"
			onscroll={handleScroll}
		>
			{#if channelStore.messageCache[channelID]?.messages.length > 0}
				<ChatSpacer {channelID} />

				{#each channelStore.messageCache[channelID].messages as message (message.id)}
					<Message server={currentServer} channel={currentChannel!} {message} />
				{/each}
			{:else}
				<ChatNoMessages channelName={currentChannel?.name} />
			{/if}
		</div>
	{/if}
</div>
