<script lang="ts">
	import type { Channel, Friend, Server } from '$lib/types/types';
	import { editorStore } from 'stores/editorStore.svelte';
	import { Editor } from '@tiptap/core';
	import { onDestroy, onMount } from 'svelte';
	import { createEditorConfig } from './richInputConfig';
	import EmojiIcon from 'ui/icons/EmojiIcon.svelte';
	import type { CreateMessageType } from '$lib/types/schemas';
	import { backend } from 'stores/backendStore.svelte';
	import { Placeholder } from '@tiptap/extensions';
	import AttachmentsButton from './attachments/AttachmentsButton.svelte';
	import Attachments from './attachments/Attachments.svelte';
	import { hasPermissions } from 'utils/permissions';
	import { page } from '$app/state';
	import { coreStore } from 'stores/coreStore.svelte';
	import RichInputList from './RichInputList.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { messageStore } from 'stores/messageStore.svelte';

	interface Props {
		server?: Server;
		channel?: Channel;
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
		const channelID = channel?.id || friend?.channel_id;
		if (!channelID) return;

		const payload: CreateMessageType = {
			server_id: server?.id || 'global',
			channel_id: channelID,
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
				onKeyupPress: () => {
					if (!channel) return;
					if (channelStore.messageCache[channel.id].scrollY > 3000) return;

					for (const message of channelStore.messageCache[channel.id].messages) {
						if (message.author.id === userStore.user?.id) {
							messageStore.editMessage = message;
							break;
						}
					}
				},
				onFocus: () => {
					editorStore.currentChannel = channel?.id || friend?.channel_id || '';
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
</script>

<div
	class="flex w-full flex-col gap-y-1 px-2.5 pb-2.5 relative"
	style="width: {coreStore.richInputLength}px;"
>
	{#if editorStore.currentInput === 'main' && editorStore.listOpen}
		<div class="absolute bottom-full left-2.5 right-2.5 mb-1 z-20">
			<RichInputList />
		</div>
	{/if}
	{#if attachments.length > 0}
		<Attachments bind:attachments />
	{/if}
	<div
		class="flex gap-x-1 box-style rounded-2xl before:rounded-[14px] focus-within:border-main-700 transition-colors duration-100"
	>
		<AttachmentsButton
			bind:attachments
			disabled={!(Boolean(friend) || hasPermissions(page.params.server_id!, 'ATTACH_FILES'))}
		/>
		<div class="relative flex w-[calc(100%-3.5rem*2)] flex-col transition duration-100">
			<div class="flex w-full">
				<div class="max-h-[10rem] w-full" bind:this={element}></div>
			</div>
		</div>
		<button
			class="h-[3.5rem] px-3 flex justify-center items-center text-main-500 hover:text-main-200 hover:cursor-pointer transition-colors duration-75 z-[1]"
		>
			<EmojiIcon />
		</button>
	</div>
</div>
