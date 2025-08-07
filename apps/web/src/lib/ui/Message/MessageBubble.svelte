<script lang="ts">
	import { generateHTML } from '@tiptap/core';
	import { messageStore } from 'stores/messageStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { MESSAGE_EXTENSIONS } from 'ui/RichInput/richInputConfig';
	import RichInputEdit from 'ui/RichInput/RichInputEdit.svelte';

	let { server, channel, message, onClickSave = $bindable() } = $props();

	// Memoize expensive HTML generation; only recompute when message.content changes
	const html = $derived.by(() => generateHTML(message.content, MESSAGE_EXTENSIONS));
</script>

<div
	class={[
		'border-[0.5px] py-1.5 px-3 w-fit max-w-full [&>*]:break-all relative z-[1]',
		message.author.id === userStore.user?.id
			? 'bg-main-900 border-main-700'
			: 'bg-main-950 border-main-800'
	]}
>
	{#if messageStore.editMessage?.id === message.id}
		<RichInputEdit {server} {channel} bind:onClickSave />
	{:else}
		{@html html}
	{/if}
</div>
