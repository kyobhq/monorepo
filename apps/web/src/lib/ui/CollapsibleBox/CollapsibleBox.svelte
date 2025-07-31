<script lang="ts">
	import type { Snippet } from 'svelte';
	import ContextMenuCategory from 'ui/ContextMenu/ContextMenuCategory.svelte';

	interface Props {
		color?: string;
		canCollapse?: boolean;
		categoryId?: string;
		header: string;
		children: Snippet;
	}

	const { children, categoryId, header, canCollapse = true, color }: Props = $props();

	let isCollapsed = $state(false);

	function collapse() {
		if (!canCollapse) return;
		isCollapsed = !isCollapsed;
	}

	function getBackgroundColor(color: string) {
		return color + '1A';
	}

	function getBorderColor(color: string) {
		return color + '80';
	}
</script>

<div
	class="border-[0.5px] border-main-400"
	style="
    background-color: {color ? getBackgroundColor(color) : 'var(--ui-main-950)'}; 
    border-color: {color ? getBorderColor(color) : 'var(--ui-main-700)'};"
>
	<button
		onclick={collapse}
		style="border-color: {color ? getBorderColor(color) : 'var(--ui-main-700)'};
      color: {color ? color : 'var(--ui-main-300)'};"
		class={[
			'relative w-full flex items-center justify-start py-1.5 px-2.5 text-main-300 text-sm',
			!isCollapsed && 'border-b-[0.5px]',
			canCollapse && 'hover:cursor-pointer hocus:bg-main-900 transition-colors duration-75'
		]}
	>
		{header}
		<ContextMenuCategory {categoryId} />
	</button>
	{#if !isCollapsed}
		<div class="p-1 flex flex-col gap-y-0.5">
			{@render children()}
		</div>
	{/if}
</div>
