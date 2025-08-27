<script lang="ts">
	import { page } from '$app/state';
	import { coreStore } from 'stores/coreStore.svelte';
	import type { Snippet } from 'svelte';
	import ContextMenuCategory from 'ui/ContextMenu/ContextMenuCategory.svelte';
	import ArrowIcon from 'ui/icons/ArrowIcon.svelte';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';
	import { hasPermissions } from 'utils/permissions';

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

	function openChannelDialog() {
		if (!categoryId) return;
		coreStore.channelDialog = {
			open: true,
			category_id: categoryId
		};
	}
</script>

<div class="overflow-hidden">
	<div class="flex items-center justify-between">
		<button
			onclick={collapse}
			class={[
				'relative w-full flex items-center justify-start py-2 px-1.5 text-main-300 text-[15px] gap-x-1.5',
				canCollapse
					? 'transition-colors duration-150 active:text-main-50/80! hover:text-main-50!'
					: 'hover:cursor-default!'
			]}
		>
			<ArrowIcon height={10} width={10} class={['transition', isCollapsed ? '-rotate-90' : '']} />
			{header}
			<ContextMenuCategory {categoryId} />
		</button>
		{#if hasPermissions(page.params.server_id!, 'MANAGE_CHANNELS')}
			<button
				class="text-main-300 hover:text-main-50! transition-colors duration-150 active:text-main-50/80!"
				onclick={openChannelDialog}
			>
				<PlusIcon height={14} width={14} />
			</button>
		{/if}
	</div>

	{#if !isCollapsed}
		<div class="flex flex-col gap-y-0.5">
			{@render children()}
		</div>
	{/if}
</div>
