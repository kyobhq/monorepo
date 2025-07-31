<script lang="ts">
	import { ContextMenu } from 'bits-ui';

	interface Props {
		onclick: () => void;
		text: string;
		destructive?: boolean;
	}

	let { onclick, text, destructive = false }: Props = $props();
</script>

<ContextMenu.Item
	class={[
		'flex px-2.5 py-1.5 hover:cursor-pointer transition-colors duration-75',
		destructive ? 'text-red-400 hocus:bg-red-400/20' : 'hocus:bg-main-800'
	]}
	{onclick}
>
	<div class="flex w-full items-center justify-between">{text}</div>
</ContextMenu.Item>
