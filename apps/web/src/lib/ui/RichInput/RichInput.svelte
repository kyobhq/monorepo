<script lang="ts">
	import type { Channel, Friend, Server } from '$lib/types/types';
	import { editorStore } from 'stores/editorStore.svelte';
	import MentionsList from './extensions/mentions/MentionsList.svelte';
	import EmojisList from './extensions/emojis/EmojisList.svelte';
	import { Editor } from '@tiptap/core';
	import { onDestroy, onMount } from 'svelte';
	import { createEditorConfig } from './richInputConfig';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';
	import EmojiIcon from 'ui/icons/EmojiIcon.svelte';

	interface Props {
		server: Server;
		channel: Channel;
		friend?: Friend;
	}

	let { server, channel, friend }: Props = $props();

	let element: Element;
	let editor: Editor;
	let attachments = $state<File[]>([]);

	async function prepareMessage(message: any) {
		if (editor.getText().length <= 0 || editor.getText().length > 2500) return;
		// const everyone = editor.getText().includes('@everyone');
		// const ids =
		// 	editor
		// 		.getText()
		// 		.match(/<@(\d+)>/g)
		// 		?.map((match) => match.slice(2, -1)) || [];

		// const payload = {
		// 	content: message,
		// 	mentions_users: [...new Set(ids)],
		// 	everyone: everyone,
		// 	attachments
		// };

		// const res = await backend.sendMessage(server.id, channel.id, payload);
		// if (res.isErr()) {
		// 	console.error(`${res.error.code}: ${res.error.error}`);
		// }
		//
		editor.commands.clearContent();
		attachments = [];
	}

	onMount(() => {
		editor = new Editor(
			createEditorConfig({
				element: element,
				placeholder: friend
					? `Message ${friend.display_name}`
					: `Message #${channel?.name} in ${server?.name}`,
				onTransaction: () => {
					editor = editor;
				},
				editorProps: {
					attributes: {
						class: 'chat-input'
					}
				},
				onEnterPress: () => prepareMessage(editor.getJSON()),
				onFocus: () => {
					editorStore.currentChannel = channel.id;
				}
			})
		);
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
		}
	});
</script>

<div class="flex w-full flex-col gap-y-1 px-2 pb-2">
	{#if editorStore.currentInput === 'main' && editorStore.currentChannel === channel.id && editorStore.mentionProps}
		<MentionsList
			props={editorStore.mentionProps}
			bind:this={editorStore.mentionsListEl}
			class="w-full"
		/>
	{/if}
	{#if editorStore.currentInput === 'main' && editorStore.currentChannel === channel.id && editorStore.emojiProps}
		<EmojisList
			props={editorStore.emojiProps}
			bind:this={editorStore.emojisListEl}
			class="w-full"
		/>
	{/if}
	<div class="flex gap-x-1">
		<button
			class="h-[3.5625rem] w-[3.5625rem] flex justify-center items-center bg-main-975 hocus:bg-main-950 border-[0.5px] border-main-800 aspect-square text-main-500 hocus:text-main-200 hover:cursor-pointer"
		>
			<PlusIcon height={22} width={22} />
		</button>
		<div
			class="bg-main-975 border-[0.5px] border-main-800 relative flex w-[calc(100%-3.5rem*2-0.5rem)] flex-col transition duration-100 focus-within:bg-main-950/70 hocus:bg-main-950"
		>
			<!-- {#if attachments.length > 0} -->
			<!-- 	<Attachments bind:attachments /> -->
			<!-- {/if} -->

			<div class="flex w-full">
				<div class="max-h-[10rem] w-full" bind:this={element}></div>
			</div>
		</div>
		<button
			class="h-[3.5625rem] w-[3.5625rem] flex justify-center items-center bg-main-975 hocus:bg-main-950 border-[0.5px] border-main-800 aspect-square text-main-500 hocus:text-main-200 hover:cursor-pointer"
		>
			<EmojiIcon height={22} width={22} />
		</button>
	</div>
</div>
