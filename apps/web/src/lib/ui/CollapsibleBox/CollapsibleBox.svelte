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
		return color + '26';
	}

	function getBorderColor(color: string) {
		return color + '59';
	}
</script>

<div
	class="box-style overflow-hidden"
	style="
    --bg-color: {color ? getBackgroundColor(color) : 'var(--ui-main-950)'};
    --border-color: {color ? getBorderColor(color) : 'var(--ui-main-850)'};
  "
>
	<button
		onclick={collapse}
		style="border-color: {color ? getBorderColor(color) : 'var(--ui-main-850)'};
      color: {color ? color : 'var(--ui-main-300)'};"
		class={[
			'relative w-full flex items-center justify-start py-1.5 px-2.5 text-main-300 text-sm rounded-t-[6px]',
			!isCollapsed && 'border-b-[2px]',
			canCollapse ? 'hover:bg-main-925 transition-colors duration-75' : 'hover:cursor-default!'
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
