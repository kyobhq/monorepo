<script lang="ts">
	import { ContextMenu } from 'bits-ui';

	interface Props {
		onclick: any;
		text: string;
		destructive?: boolean;
		success?: boolean;
	}

	let { onclick, text, success = false, destructive = false }: Props = $props();
</script>

<ContextMenu.Item
	class={[
		'flex px-2.5 py-1.5 hover:cursor-pointer transition-colors duration-75 select-none text-sm rounded-lg',
		!destructive && !success && 'hover:bg-main-850',
		destructive && 'text-red-400 hover:bg-red-400/20',
		success && 'text-green-400 hover:bg-green-400/20'
	]}
	{onclick}
>
	<div class="flex w-full items-center justify-between">{text}</div>
</ContextMenu.Item>
