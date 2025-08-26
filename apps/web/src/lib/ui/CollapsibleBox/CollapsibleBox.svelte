<script lang="ts">
	import type { Snippet } from 'svelte';
	import ContextMenuCategory from 'ui/ContextMenu/ContextMenuCategory.svelte';
	import ArrowIcon from 'ui/icons/ArrowIcon.svelte';

	interface Props {
		canCollapse?: boolean;
		categoryId?: string;
		header: string;
		children: Snippet;
	}

	const { children, categoryId, header, canCollapse = true }: Props = $props();

	let isCollapsed = $state(false);

	function collapse() {
		if (!canCollapse) return;
		isCollapsed = !isCollapsed;
	}
</script>

<div class="overflow-hidden">
	<button
		onclick={collapse}
		class={[
			'relative w-full flex items-center justify-start py-1.5 px-1.5 text-main-300 text-sm gap-x-1.5',
			canCollapse
				? 'transition-colors duration-150 active:text-main-50/80! hover:text-main-50!'
				: 'hover:cursor-default!'
		]}
	>
		{header}
		<ArrowIcon height={10} width={10} class={['transition', isCollapsed ? '-rotate-90' : '']} />
		<ContextMenuCategory {categoryId} />
	</button>
	{#if !isCollapsed}
		<div class="flex flex-col gap-y-0.5">
			{@render children()}
		</div>
	{/if}
</div>
