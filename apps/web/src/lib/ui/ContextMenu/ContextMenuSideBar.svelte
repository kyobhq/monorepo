<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import ContextMenuSkeleton from './ContextMenuSkeleton.svelte';
	import ContextMenuItem from './ContextMenuItem.svelte';
	import { hasPermissions } from 'utils/permissions';
	import { page } from '$app/state';
	import FilledPlusIcon from 'ui/icons/FilledPlusIcon.svelte';
</script>

{#if hasPermissions(page.params.server_id!, 'MANAGE_CHANNELS')}
	<ContextMenuSkeleton>
		{#snippet contextMenuContent()}
			<ContextMenuItem
				Icon={FilledPlusIcon}
				onclick={() => (coreStore.categoryDialog = true)}
				text="Create category"
			/>
		{/snippet}
	</ContextMenuSkeleton>
{/if}
