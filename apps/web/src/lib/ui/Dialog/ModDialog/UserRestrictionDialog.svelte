<script lang="ts">
	import ModDialog from './ModDialog.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
</script>

<ModDialog
	title={coreStore.restrictionDialog.title}
	open={coreStore.restrictionDialog.open}
	bind:openState={coreStore.restrictionDialog.open}
>
	<div class="flex flex-col mt-4">
		<p class="text-sm text-main-500">Reason of your {coreStore.restrictionDialog.restriction}:</p>
		<p class="text-lg">{coreStore.restrictionDialog.reason}</p>
	</div>
</ModDialog>
