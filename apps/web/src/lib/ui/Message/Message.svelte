<script lang="ts">
	import type { Channel, Friend, Message, Server } from '$lib/types/types';
	import { messageStore } from 'stores/messageStore.svelte';
	import ContextMenuMessage from 'ui/ContextMenu/ContextMenuMessage.svelte';
	import MessageEditSentence from './MessageEditSentence.svelte';
	import MessageBubble from './MessageBubble.svelte';
	import MessageAuthor from './MessageAuthor.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import MessageAvatar from './MessageAvatar.svelte';

	interface Props {
		server?: Server;
		channel?: Channel;
		friend?: Friend;
		message: Message;
	}

	let { server, channel, friend, message }: Props = $props();

	const messageIsRecent = $derived(channelStore.messageIsRecent(channel?.id || '', message.id));

	let onClickSave = $state<() => Promise<void>>();
	let author = $derived(messageStore.getAuthor(message.author.id!) || message.author);
	let hover = $state(false);
</script>

<div
	role="group"
	class={[
		'flex gap-x-2.5 items-end hover:bg-main-950/50 px-6 transition-colors duration-75 relative first:mb-0',
		messageIsRecent ? 'mb-0.5' : 'mb-4'
	]}
	onmouseenter={() => (hover = true)}
	onmouseleave={() => (hover = false)}
>
	{#if !messageIsRecent}
		<MessageAvatar {author} hoverAvatar={hover} />
	{/if}
	<div class="flex flex-col gap-y-1 relative w-[calc(100%-4rem)]">
		{#if messageStore.editMessage?.id === message.id}
			<MessageEditSentence bind:onClickSave />
		{/if}
		<MessageBubble
			{server}
			{channel}
			{friend}
			{message}
			{messageIsRecent}
			hoverMessage={hover}
			bind:onClickSave
		/>
		{#if !messageIsRecent}
			<MessageAuthor {author} {message} />
		{/if}
	</div>

	<ContextMenuMessage {message} />
</div>
