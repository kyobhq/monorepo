<script lang="ts">
	import type { Channel, Friend, Message, Server } from '$lib/types/types';
	import { messageStore } from 'stores/messageStore.svelte';
	import ContextMenuMessage from 'ui/ContextMenu/ContextMenuMessage.svelte';
	import MessageEditSentence from './MessageEditSentence.svelte';
	import MessageAuthor from './MessageAuthor.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import MessageAvatar from './MessageAvatar.svelte';
	import MessageContent from './MessageContent.svelte';
	import { userStore } from 'stores/userStore.svelte';

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
	let mentioned = $derived(userStore.user && message.mentions_users?.includes(userStore.user.id));
	let hover = $state(false);
</script>

<div
	role="group"
	class={[
		'flex gap-x-2.5 items-start px-6 transition-colors duration-75 relative first:mb-0',
		!messageIsRecent && 'mt-4 py-1',
		mentioned
			? 'bg-mention/20 mention-bar hover:bg-mention/35 mix-blend-plus-lighter'
			: 'hover:bg-main-950/50'
	]}
	onmouseenter={() => (hover = true)}
	onmouseleave={() => (hover = false)}
>
	{#if !messageIsRecent}
		<MessageAvatar {author} hoverAvatar={hover} />
	{/if}
	<div class="flex flex-col relative w-[calc(100%-4rem)]">
		{#if !messageIsRecent}
			<MessageAuthor {author} {message} />
		{/if}
		<MessageContent
			{server}
			{channel}
			{friend}
			{message}
			{messageIsRecent}
			hoverMessage={hover}
			bind:onClickSave
		/>
		{#if messageStore.editMessage?.id === message.id}
			<MessageEditSentence bind:onClickSave {messageIsRecent} />
		{/if}
	</div>

	<ContextMenuMessage {message} />
</div>
