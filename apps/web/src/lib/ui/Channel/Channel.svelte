<script lang="ts">
	import VolumeIcon from '../icons/VolumeIcon.svelte';
	import Lock from '../icons/Lock.svelte';
	import HashChat from '../icons/HashChat.svelte';
	import type { ChannelTypes } from '$lib/types/types';

	interface Props {
		type: ChannelTypes;
		name: string;
		serverName?: string;
		onclick: () => void;
		active?: boolean;
	}

	const { type, name, serverName, onclick, active }: Props = $props();

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
		'flex items-center w-full gap-x-2 hover:cursor-pointer transition-colors duration-75 py-2 px-2.5',
		active ? 'text-main-50 bg-main-900' : 'hocus:bg-main-900 hocus:text-main-200 text-main-300'
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
</button>
