<script lang="ts">
	import type { Message } from '$lib/types/types';
	import { formatMessageTime } from 'utils/time';

	interface Props {
		message: Message;
	}

	let { message }: Props = $props();
</script>

<div class="flex items-baseline gap-x-1.5 select-none w-fit relative z-[1]">
	<p>{message.author.display_name}</p>
	<time class="text-sm text-main-600">{formatMessageTime(message.created_at)}</time>
	{#if new Date(message.created_at) < new Date(message.updated_at)}
		<p class="text-sm text-main-600">(edited)</p>
	{/if}
</div>
