<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import type { Component } from 'svelte';

	interface Props {
		onclick: any;
		text: string;
		destructive?: boolean;
		success?: boolean;
		Icon: Component;
	}

	let { onclick, text, success = false, destructive = false, Icon }: Props = $props();
</script>

<ContextMenu.Item
	class={[
		'flex items-center gap-x-3 px-2.5 py-2 hover:cursor-pointer transition-colors duration-75 select-none text-sm rounded-lg',
		!destructive && !success && 'hover:bg-main-850',
		destructive && 'text-red-400 hover:bg-red-400/20',
		success && 'text-green-400 hover:bg-green-400/20'
	]}
	{onclick}
>
	<Icon height={20} width={20} />
	<div class="flex w-full items-center justify-between">{text}</div>
</ContextMenu.Item>
