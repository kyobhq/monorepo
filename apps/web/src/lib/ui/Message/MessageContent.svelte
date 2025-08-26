<script lang="ts">
	import type { Message, Channel, Server, Friend } from '$lib/types/types';
	import { generateHTML } from '@tiptap/core';
	import { coreStore } from 'stores/coreStore.svelte';
	import { messageStore } from 'stores/messageStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { onMount } from 'svelte';
	import { MESSAGE_EXTENSIONS } from 'ui/RichInput/richInputConfig';
	import RichInputEdit from 'ui/RichInput/RichInputEdit.svelte';
	import MessageAttachments from './MessageAttachments.svelte';
	import { generateTextWithExt } from 'utils/richInput';

	interface Props {
		server?: Server;
		channel?: Channel;
		friend?: Friend;
		hoverMessage: boolean;
		message: Message;
		messageIsRecent: boolean;
		onClickSave?: () => Promise<void>;
	}

	let {
		server,
		channel,
		friend,
		message,
		messageIsRecent,
		hoverMessage,
		onClickSave = $bindable()
	}: Props = $props();

	let messageEl = $state<HTMLElement>();
	const html = $derived.by(() => generateHTML(message.content, MESSAGE_EXTENSIONS));
	const messageLength = $derived.by(() => generateTextWithExt(message.content).length);

	function handleMention(e: MouseEvent) {
		const target = e.target as HTMLButtonElement;
		const userId = target.attributes.getNamedItem('user-id')?.value;
		if (userId) {
			coreStore.openProfile(userId, target, 'right');
		}
	}

	onMount(() => {
		const mentions = messageEl?.querySelectorAll('[data-type="mention"]');

		if (!mentions) return;
		for (const mention of mentions) {
			(mention as HTMLButtonElement).addEventListener('click', handleMention);
		}

		return () => {
			for (const mention of mentions) {
				(mention as HTMLButtonElement).removeEventListener('click', handleMention);
			}
		};
	});
</script>

<div class={['flex flex-col gap-y-1.5', messageIsRecent && 'pl-[3.625rem]']}>
	{#if messageLength > 0}
		{#if messageStore.editMessage?.id === message.id}
			<RichInputEdit {server} {channel} {friend} bind:onClickSave />
		{:else}
			<div
				bind:this={messageEl}
				class={[
					'py-0.5 w-fit max-w-full [&>*]:break-all relative z-[1] rounded-2xl transition-all',
					!messageIsRecent && ' rounded-bl-sm'
				]}
			>
				{@html html}
			</div>
		{/if}
	{/if}
	{#if message.attachments?.length > 0}
		<MessageAttachments attachments={message.attachments} hover={hoverMessage} />
	{/if}
</div>
