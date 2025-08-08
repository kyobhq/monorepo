<script lang="ts">
	import type { Message, Channel, Server } from '$lib/types/types';
	import { generateHTML } from '@tiptap/core';
	import { coreStore } from 'stores/coreStore.svelte';
	import { messageStore } from 'stores/messageStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { onMount } from 'svelte';
	import { MESSAGE_EXTENSIONS } from 'ui/RichInput/richInputConfig';
	import RichInputEdit from 'ui/RichInput/RichInputEdit.svelte';

	interface Props {
		server: Server;
		channel: Channel;
		message: Message;
		onClickSave?: () => Promise<void>;
	}

	let { server, channel, message, onClickSave = $bindable() }: Props = $props();

	let messageEl = $state<HTMLElement>();
	const html = $derived.by(() => generateHTML(message.content, MESSAGE_EXTENSIONS));
	let mentioned = $derived(userStore.user && message.mentions_users?.includes(userStore.user.id));

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

<div
	bind:this={messageEl}
	class={[
		'border-[0.5px] py-1.5 px-3 w-fit max-w-full [&>*]:break-all relative z-[1]',
		message.author.id === userStore.user?.id
			? 'bg-main-900 border-main-700'
			: 'bg-main-950 border-main-800',
		mentioned && 'border-mention! text-mention! bg-mention/20!'
	]}
>
	{#if messageStore.editMessage?.id === message.id}
		<RichInputEdit {server} {channel} bind:onClickSave />
	{:else}
		{@html html}
	{/if}
</div>
