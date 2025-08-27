<script lang="ts">
	import { page } from '$app/state';
	import type { Member, Message } from '$lib/types/types';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import ContextMenuUser from 'ui/ContextMenu/ContextMenuUser.svelte';
	import { formatMessageTime } from 'utils/time';

	interface Props {
		author: Partial<Member>;
		message: Message;
	}

	let { author, message }: Props = $props();
	let displayNameEl = $state<HTMLButtonElement>();
	let topRole = $derived(serverStore.getUserTopRole(page.params.server_id!, author.id!));
</script>

<div class="flex items-baseline gap-x-1.5 select-none w-fit relative z-[1]">
	<button
		bind:this={displayNameEl}
		class="relative hover:underline hover:cursor-pointer"
		style="color: {topRole ? topRole.color : 'var(--ui-main-50)'}"
		onclick={() => {
			if (author.id === userStore.user!.id) {
				coreStore.openMyProfile(displayNameEl!, 'right');
			} else {
				coreStore.openProfile(author.id!, displayNameEl!, 'right');
			}
		}}
	>
		{author.display_name}
		<ContextMenuUser memberID={author.id!} />
	</button>
	<time class="text-xs text-main-600">{formatMessageTime(message.created_at)}</time>
	{#if new Date(message.created_at) < new Date(message.updated_at)}
		<p class="text-xs text-main-600">[edited]</p>
	{/if}
</div>
