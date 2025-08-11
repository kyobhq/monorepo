<script lang="ts">
	import type { Channel, Server } from '$lib/types/types';
	import { editorStore } from 'stores/editorStore.svelte';
	import MentionsList from './extensions/mentions/MentionsList.svelte';
	import EmojisList from './extensions/emojis/EmojisList.svelte';
	import { Editor } from '@tiptap/core';
	import { onDestroy, onMount } from 'svelte';
	import { createEditorConfig } from './richInputConfig';
	import { Placeholder } from '@tiptap/extensions';
	import { messageStore } from 'stores/messageStore.svelte';
	import type { EditMessageType } from '$lib/types/schemas';
	import { backend } from 'stores/backendStore.svelte';
	import { page } from '$app/state';
	import { coreStore } from 'stores/coreStore.svelte';

	interface Props {
		server: Server;
		channel: Channel;
		onClickSave?: () => Promise<void>;
	}

	let { server, channel, onClickSave = $bindable() }: Props = $props();

	let element: Element;
	let editor: Editor;

	async function prepareMessage(message: any) {
		if (!messageStore.editMessage) return;
		const messageId = messageStore.editMessage.id;

		if (editor.getText().trim().length <= 0) {
			coreStore.destructiveDialog = {
				open: true,
				title: `Delete Message`,
				subtitle: 'Are you sure you want to delete this message?',
				buttonText: 'Delete',
				onclick: async () => {
					coreStore.destructiveDialog.open = false;
					const res = await backend.deleteMessage(messageId, {
						server_id: page.params.server_id!,
						channel_id: page.params.channel_id!,
						author_id: messageStore.editMessage!.author.id
					});

					res.match(
						() => {},
						(error) => {
							console.error(`${error.code}: ${error.message}`);
						}
					);
				}
			};
			return;
		}

		if (editor.getText().length > 2500) return;
		const everyone = editor.getText().includes('@everyone');
		const ids =
			editor
				.getText()
				.match(/<@(\d+)>/g)
				?.map((match) => match.slice(2, -1)) || [];

		const payload: EditMessageType = {
			server_id: server.id,
			channel_id: channel.id,
			content: message,
			mentions_users: [...new Set(ids)],
			mentions_roles: [],
			mentions_channels: [],
			everyone: everyone
		};

		messageStore.stopEditing();
		const res = await backend.editMessage(messageId, payload);
		res.match(
			() => {},
			(error) => {
				console.error(`${error.code}: ${error.message}`);
			}
		);
	}

	onMount(() => {
		editorStore.currentInput = 'edit';
		editor = new Editor(
			createEditorConfig({
				element: element,
				autofocus: 'end',
				content: messageStore.editMessage?.content,
				placeholder: Placeholder.configure({
					placeholder: ''
				}),
				onTransaction: () => {
					editor = editor;
				},
				editorProps: {
					attributes: {
						class: 'edit-chat-input'
					}
				},
				onEnterPress: () => prepareMessage(editor.getJSON()),
				onEscapePress: () => messageStore.stopEditing()
			})
		);

		onClickSave = () => prepareMessage(editor.getJSON());
	});

	onDestroy(() => {
		editorStore.currentInput = 'main';
		if (editor) {
			editor.destroy();
		}
	});

	$effect(() => {
		if (!editor) return;
		editor.setOptions();
	});
</script>

<div class="flex w-full flex-col gap-y-1">
	{#if editorStore.currentInput === 'edit' && editorStore.mentionProps}
		<MentionsList
			props={editorStore.mentionProps}
			bind:this={editorStore.mentionsListEl}
			class="w-full"
		/>
	{/if}
	{#if editorStore.currentInput === 'edit' && editorStore.emojiProps}
		<EmojisList
			props={editorStore.emojiProps}
			bind:this={editorStore.emojisListEl}
			class="w-full"
		/>
	{/if}
	<div class="max-h-[10rem] w-full" bind:this={element}></div>
</div>
