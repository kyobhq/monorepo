<script lang="ts">
	import VolumeIcon from '../icons/VolumeIcon.svelte';
	import Lock from '../icons/Lock.svelte';
	import HashChat from '../icons/HashChat.svelte';
	import type { ChannelTypes } from '$lib/types/types';
	import ContextMenuChannel from 'ui/ContextMenu/ContextMenuChannel.svelte';

	interface Props {
		id: string;
		categoryId: string;
		type: ChannelTypes;
		name: string;
		serverName?: string;
		onclick: () => void;
		active?: boolean;
		unread: boolean;
		mentions: number;
	}

	const {
		id,
		categoryId,
		type,
		name,
		serverName,
		onclick,
		active,
		unread = false,
		mentions
	}: Props = $props();

	const ICONS = {
		textual: HashChat,
		voice: VolumeIcon,
		'textual-e2ee': Lock,
		dm: Lock
	};
	const Icon = ICONS[type];
</script>

<button
	{onclick}
	class={[
		'relative flex items-center w-full gap-x-2 hover:cursor-pointer transition duration-150 py-2 px-2.5 rounded-[10px] active-scale-down',
		active ? 'text-main-50 bg-main-850' : 'hover:bg-main-900 hover:text-main-50 text-main-300',
		unread && 'text-main-50!'
	]}
>
	<Icon height={18} width={18} />
	<div class="flex items-baseline gap-x-1.5">
		{name}
		{#if serverName}
			<span class="text-main-600 text-xs">
				in <span class="font-semibold">{serverName}</span>
			</span>
		{/if}
	</div>

	{#if mentions > 0}
		<div
			class="h-5 w-5 absolute right-3 top-1/2 -translate-y-1/2 bg-red-400 rounded-md flex items-center justify-center text-sm font-semibold"
		>
			{mentions}
		</div>
	{:else if unread}
		<div class="h-2 w-2 absolute right-3 top-1/2 -translate-y-1/2 bg-white rounded-full"></div>
	{/if}

	<ContextMenuChannel {categoryId} channelId={id} />
</button>
