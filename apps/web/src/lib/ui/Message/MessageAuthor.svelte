<script lang="ts">
	import type { Message, User } from '$lib/types/types';
	import { coreStore } from 'stores/coreStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { formatMessageTime } from 'utils/time';

	interface Props {
		author: Partial<User>;
		message: Message;
	}

	let { author, message }: Props = $props();
	let displayNameEl = $state<HTMLButtonElement>();
</script>

<div class="flex items-baseline gap-x-1.5 select-none w-fit relative z-[1]">
	<button
		bind:this={displayNameEl}
		class="hover:underline hover:cursor-pointer"
		onclick={() => {
			if (author.id === userStore.user!.id) {
				coreStore.openMyProfile(displayNameEl!, 'right');
			}
		}}>{author.display_name}</button
	>
	<time class="text-xs text-main-600">{formatMessageTime(message.created_at)}</time>
	{#if new Date(message.created_at) < new Date(message.updated_at)}
		<p class="text-xs text-main-600">[edited]</p>
	{/if}
</div>
