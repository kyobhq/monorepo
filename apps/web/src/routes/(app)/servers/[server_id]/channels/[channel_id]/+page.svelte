<script lang="ts">
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { tick } from 'svelte';
	import ChannelHeader from 'ui/ChannelHeader/ChannelHeader.svelte';
	import Message from 'ui/Message/Message.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';

	let scrollContainer = $state<HTMLDivElement>();
	let messageCount = $state(0);
	let isAtBottom = $state(true);

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

	async function getMessages(channelID: string) {
		const res = await backend.getMessages(channelID);
		res.match(
			(messages) => {
				channelStore.messages = messages ? messages.reverse() : [];
				tick().then(() => scrollToBottom(false));
			},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}

	$effect(() => {
		if (page.params.channel_id) getMessages(page.params.channel_id);
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
		class="flex flex-col-reverse w-full h-full gap-y-4 overflow-auto pb-4 pt-18"
		bind:this={scrollContainer}
		onscroll={handleScroll}
	>
		{#if channelStore.messages.length > 0}
			{#each channelStore.messages as message (message.id)}
				<Message server={currentServer} channel={currentChannel} {message} />
			{/each}
		{:else}
			No messages in {currentChannel?.name} yet :c
		{/if}
	</div>

	<RichInput server={currentServer} channel={currentChannel} />
{/if}
