<script lang="ts">
	import { beforeNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import { backend } from 'stores/backendStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { onMount, tick } from 'svelte';
	import StraightFaceEmoji from 'ui/icons/StraightFaceEmoji.svelte';
	import Message from 'ui/Message/Message.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';

	let scrollContainer = $state<HTMLDivElement>();
	let messageCount = $state(0);
	let isAtBottom = $state(true);

	const currentFriend = $derived(
		userStore.friends.find((friend) => friend.channel_id === page.params.channel_id)
	);

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
		const res = await backend.getMessages('global', channelID);
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

	onMount(async () => {
		if (page.params.channel_id) await getMessages(page.params.channel_id);
	});

	beforeNavigate(async ({ from, to }) => {
		if (
			from?.params?.channel_id &&
			to?.params?.channel_id &&
			from?.params?.channel_id !== to?.params?.channel_id
		) {
			await getMessages(to.params.channel_id);
		}
	});

	$effect(() => {
		if (channelStore.messages.length !== messageCount && isAtBottom) {
			tick().then(() => scrollToBottom(false));
		}
		messageCount = channelStore.messages.length;
	});
</script>

{#if currentFriend}
	<div
		class="flex flex-col-reverse w-full h-full gap-y-4 overflow-auto pb-4 pt-18"
		bind:this={scrollContainer}
		onscroll={handleScroll}
	>
		{#if channelStore.messages.length > 0}
			{#each channelStore.messages as message (message.id)}
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
