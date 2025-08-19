<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { onMount, tick } from 'svelte';
	import ChannelHeader from 'ui/ChannelHeader/ChannelHeader.svelte';
	import StraightFaceEmoji from 'ui/icons/StraightFaceEmoji.svelte';
	import Message from 'ui/Message/Message.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';

	let scrollContainer = $state<HTMLDivElement>();
	let messageCount = $state(0);
	let isAtBottom = $state(true);
	let messagesLoaded = $state(false);

	const currentServer = $derived.by(() => {
		if (!page.params.server_id) return;
		return serverStore.getServer(page.params.server_id);
	});
	const currentChannel = $derived.by(() => {
		if (!page.params.server_id || !page.params.channel_id) return;
		return channelStore.getChannel(page.params.server_id, page.params.channel_id);
	});

	function handleScroll(ev: Event) {
		// 0 is bottom
		isAtBottom = Math.abs((ev.target as HTMLDivElement).scrollTop) <= 100;
	}

	const scrollToBottom = (smooth = false) => {
		if (scrollContainer) {
			scrollContainer.scrollTo({
				top: scrollContainer.scrollHeight,
				behavior: smooth ? 'smooth' : 'instant'
			});
		}
	};

	async function getMessages(serverID: string, channelID: string) {
		const res = await backend.getMessages(serverID, channelID);
		res.match(
			(messages) => {
				channelStore.messages = messages ? messages.reverse() : [];
				tick().then(() => scrollToBottom(false));
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

	$effect(() => {
		if (channelStore.messages.length !== messageCount && isAtBottom) {
			tick().then(() => scrollToBottom(false));
		}
		messageCount = channelStore.messages.length;
	});
</script>

{#if currentChannel && currentServer}
	<ChannelHeader name={currentChannel.name} description={currentChannel.description} />

	<div
		class="flex flex-col-reverse w-full h-full overflow-auto pb-4 pt-18"
		bind:this={scrollContainer}
		onscroll={handleScroll}
	>
		{#if messagesLoaded}
			{#if channelStore.messages.length > 0}
				{#each channelStore.messages as message (message.id)}
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
