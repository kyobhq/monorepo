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
	}

	const { id, categoryId, type, name, serverName, onclick, active }: Props = $props();

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
		'relative flex items-center w-full gap-x-2 hover:cursor-pointer transition duration-150 py-2 px-2.5 rounded-md active-scale-down',
		active ? 'text-main-50 active-channel' : 'hover:bg-main-900 hover:text-main-200 text-main-300'
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
	<ContextMenuChannel {categoryId} channelId={id} />
</button>
