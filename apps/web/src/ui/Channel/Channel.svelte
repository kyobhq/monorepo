<script lang="ts">
	import VolumeIcon from 'ui/icons/VolumeIcon.svelte';
	import Lock from 'ui/icons/Lock.svelte';
	import HashChat from 'ui/icons/HashChat.svelte';

	interface Props {
		type: 'textual' | 'voice' | 'textual-e2ee';
		title: string;
		subtitle?: string;
		onclick: () => void;
		active?: boolean;
	}

	const { type, title, subtitle, onclick, active }: Props = $props();

	const ICONS = {
		textual: HashChat,
		voice: VolumeIcon,
		'textual-e2ee': Lock
	};
	const Icon = ICONS[type];
</script>

<button
	{onclick}
	class={[
		'flex items-center w-full gap-x-2 text-sm hover:cursor-pointer transition-colors duration-75 py-2 px-2.5',
		active ? 'text-main-50 bg-main-900' : 'hocus:bg-main-900 hocus:text-main-200 text-main-300'
	]}
>
	<Icon height={18} width={18} />
	{title}
	{#if subtitle}
		<span class="text-main-600 text-xs">
			in <span class="font-semibold">{subtitle}</span>
		</span>
	{/if}
</button>
