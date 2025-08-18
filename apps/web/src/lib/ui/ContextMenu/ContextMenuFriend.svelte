<script lang="ts">
	import ContextMenuSkeleton from './ContextMenuSkeleton.svelte';
	import ContextMenuItem from './ContextMenuItem.svelte';
	import type { Friend } from '$lib/types/types';

	interface Props {
		friend: Friend;
	}

	let { friend }: Props = $props();
</script>

<ContextMenuSkeleton>
	{#snippet contextMenuContent()}
		<ContextMenuItem onclick={() => {}} text="Call" />
		<ContextMenuItem onclick={() => {}} text="Remove Friend" destructive />
		<ContextMenuItem onclick={() => {}} text="Block" destructive />
	{/snippet}
</ContextMenuSkeleton>
