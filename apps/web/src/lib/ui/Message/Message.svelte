<script lang="ts">
	import type { Channel, Message, Server } from '$lib/types/types';
	import { messageStore } from 'stores/messageStore.svelte';
	import ContextMenuMessage from 'ui/ContextMenu/ContextMenuMessage.svelte';
	import MessageEditSentence from './MessageEditSentence.svelte';
	import MessageBubble from './MessageBubble.svelte';
	import MessageAuthor from './MessageAuthor.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';

	interface Props {
		server: Server;
		channel: Channel;
		message: Message;
	}

	let { server, channel, message }: Props = $props();

	let avatarEl = $state<HTMLButtonElement>();
	let onClickSave = $state<() => Promise<void>>();
	let author = $derived(
		message.author.id === userStore.user!.id ? userStore.user! : message.author
	);
</script>

<div
	class="flex gap-x-2.5 items-end hover:bg-main-950/50 px-6 py-1.5 transition-colors duration-75 relative"
>
	<button
		bind:this={avatarEl}
		class="h-12 w-12 relative highlight-border mb-1 select-none shrink-0 hover:after:border-main-50/75 hover:cursor-pointer rounded-md overflow-hidden"
		onclick={() => {
			if (author.id === userStore.user!.id) {
				coreStore.openMyProfile(avatarEl!, 'right');
			} else {
				coreStore.openProfile(author.id, avatarEl!, 'right');
			}
		}}
	>
		<img src={author.avatar} alt="" class="w-full h-full object-cover" />
	</button>
	<div class="flex flex-col gap-y-1 relative w-[calc(100%-4rem)]">
		{#if messageStore.editMessage?.id === message.id}
			<MessageEditSentence bind:onClickSave />
		{/if}
		<MessageBubble {server} {channel} {message} bind:onClickSave />
		<MessageAuthor {author} {message} />
	</div>

	<ContextMenuMessage {message} />
</div>
