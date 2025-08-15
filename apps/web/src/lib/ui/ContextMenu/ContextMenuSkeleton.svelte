<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import type { Snippet } from 'svelte';

	interface Props {
		contextMenuContent: Snippet;
	}

	let { contextMenuContent }: Props = $props();

	function contextMenuMouseDown(e: MouseEvent) {
		e.stopImmediatePropagation();
	}
</script>

<ContextMenu.Root>
	<ContextMenu.Trigger class="absolute inset-0"></ContextMenu.Trigger>
	<ContextMenu.Portal>
		<ContextMenu.Content
			class="bg-main-950 border-[2px] border-main-850 relative flex w-[225px] flex-col gap-y-1 p-1 select-none z-10 rounded-xl shadow-box"
			onmousedown={contextMenuMouseDown}
		>
			{@render contextMenuContent()}
		</ContextMenu.Content>
	</ContextMenu.Portal>
</ContextMenu.Root>
