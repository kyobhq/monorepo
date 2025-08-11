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
	import type { CreateMessageType } from '$lib/types/schemas';
	import { backend } from 'stores/backendStore.svelte';
	import { Placeholder } from '@tiptap/extensions';
	import AttachmentsButton from './attachments/AttachmentsButton.svelte';
	import Attachments from './attachments/Attachments.svelte';

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
		if (editor.getText().length > 2500) return;
		if (editor.getText().trim().length <= 0 && attachments.length <= 0) return;
		const everyone = editor.getText().includes('@everyone');
		const ids =
			editor
				.getText()
				.match(/<@(.*)>/g)
				?.map((match) => match.slice(2, -1)) || [];

		const payload: CreateMessageType = {
			server_id: server.id,
			channel_id: channel.id,
			content: message,
			mentions_users: [...new Set(ids)],
			mentions_roles: [],
			mentions_channels: [],
			everyone: everyone,
			attachments
		};

		const res = await backend.createMessage(payload);
		res.match(
			() => {},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);

		editor.commands.clearContent();
		attachments = [];
	}

	onMount(() => {
		editor = new Editor(
			createEditorConfig({
				element: element,
				placeholder: Placeholder.configure({
					placeholder: () =>
						friend
							? `Message ${friend.display_name}`
							: `Message #${channel?.name} in ${server?.name}`
				}),
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

	$effect(() => {
		if (!editor) return;
		editor.setOptions();
	});

	const icon =
		'h-[3.5625rem] w-[3.5625rem] flex justify-center items-center bg-main-900 hocus:bg-main-800/80 border-[0.5px] border-main-700 aspect-square text-main-500 hocus:text-main-200 hover:cursor-pointer transition-colors duration-75';
</script>

<div class="flex w-full flex-col gap-y-1 px-2.5 pb-2.5">
	{#if editorStore.currentInput === 'main' && editorStore.mentionProps}
		<MentionsList
			props={editorStore.mentionProps}
			bind:this={editorStore.mentionsListEl}
			class="w-full"
		/>
	{/if}
	{#if editorStore.currentInput === 'main' && editorStore.emojiProps}
		<EmojisList
			props={editorStore.emojiProps}
			bind:this={editorStore.emojisListEl}
			class="w-full"
		/>
	{/if}
	{#if attachments.length > 0}
		<Attachments bind:attachments />
	{/if}
	<div class="flex gap-x-1">
		<AttachmentsButton bind:attachments />
		<div
			class="bg-main-900 border-[0.5px] border-main-700 relative flex w-[calc(100%-3.5rem*2-0.625rem)] flex-col transition duration-100 focus-within:border-main-500 hocus:bg-main-800/70 rounded-[2px]"
		>
			<div class="flex w-full">
				<div class="max-h-[10rem] w-full" bind:this={element}></div>
			</div>
		</div>
		<button class={[icon, 'rounded-l-[2px] rounded-r-md']}>
			<EmojiIcon height={22} width={22} />
		</button>
	</div>
</div>
