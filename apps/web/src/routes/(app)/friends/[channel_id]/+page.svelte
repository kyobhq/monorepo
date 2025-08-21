<script lang="ts">
	import { afterNavigate, beforeNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import { channelStore } from 'stores/channelStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { onMount } from 'svelte';
	import StraightFaceEmoji from 'ui/icons/StraightFaceEmoji.svelte';
	import Message from 'ui/Message/Message.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';

	const SCROLL_THRESHOLD = 500;
	const currentFriend = $derived(
		userStore.friends.find((friend) => friend.channel_id === page.params.channel_id)
	);

	let messagesLoaded = $state(false);
	let canLoadMore = $state(true);
	let showMessages = $derived(
		coreStore.firstLoad.sidebar && coreStore.firstLoad.serverbar && messagesLoaded
	);

	async function handleScroll(e: Event) {
		if (!page.params.channel_id || !canLoadMore) return;
		const target = e.target as HTMLDivElement;
		const scrollY = target.clientHeight + Math.abs(target.scrollTop);

		if (target.scrollHeight - SCROLL_THRESHOLD <= scrollY) {
			canLoadMore = false;
			const hasMore = await channelStore.loadMoreMessages(
				'global',
				page.params.channel_id,
				'before'
			);

			if (hasMore) {
				setTimeout(() => {
					canLoadMore = true;
				}, 500);
			}
		} else if (
			channelStore.messageCache[page.params.channel_id]?.offsetHeight > 0 &&
			scrollY -
				channelStore.messageCache[page.params.channel_id]?.offsetHeight -
				target.clientHeight <=
				SCROLL_THRESHOLD
		) {
			canLoadMore = false;

			const hasMore = await channelStore.loadMoreMessages(
				'global',
				page.params.channel_id,
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
		if (page.params.channel_id) {
			await channelStore.ensureMessagesLoaded('global', page.params.channel_id);
			messagesLoaded = true;
		}
	});

	afterNavigate(async ({ from, to }) => {
		const fromChannelID = from?.params?.channel_id;
		const toChannelID = to?.params?.channel_id;
		if (!fromChannelID || !toChannelID) return;

		const sameChannel = fromChannelID === toChannelID;
		if (!sameChannel) {
			await channelStore.ensureMessagesLoaded('global', toChannelID);
			messagesLoaded = true;
			canLoadMore = !(channelStore.messageCache[toChannelID]?.hasReachedEnd ?? false);
		}
	});
</script>

{#if currentFriend}
	<div
		class="flex flex-col-reverse w-full h-full gap-y-4 overflow-auto pb-4 pt-18"
		onscroll={handleScroll}
	>
		{#if showMessages && channelStore.messageCache[currentFriend.channel_id!]?.messages.length > 0}
			{#each channelStore.messageCache[currentFriend.channel_id!].messages as message (message.id)}
				<Message friend={currentFriend} {message} />
			{/each}
		{:else}
			<div
				class="h-full w-full flex items-center justify-center text-2xl flex-col gap-y-4 text-main-700"
			>
				<StraightFaceEmoji height={128} width={128} />
				<p>No messages with <span class="font-medium">{currentFriend.display_name}</span> yet</p>
			</div>
		{/if}
	</div>

	<RichInput friend={currentFriend} />
{/if}
