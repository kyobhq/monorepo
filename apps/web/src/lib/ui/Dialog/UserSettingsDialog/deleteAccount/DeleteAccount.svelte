<script lang="ts">
	import { goto } from '$app/navigation';
	import { backend } from 'stores/backendStore.svelte';
	import DestructiveBar from 'ui/DestructiveBar/DestructiveBar.svelte';
	import { logErr } from 'utils/print';

	let isSubmitting = $state(false);
	let activateDelete = $state(false);

	async function handleDelete() {
		isSubmitting = true;
		const res = await backend.deleteAccount();
		res.match(
			() => goto('/signin'),
			(err) => logErr(err)
		);
		isSubmitting = false;
	}
</script>

<div>
	<p class="font-medium select-none">Account removal</p>
	<p class="text-sm text-main-500 select-none">
		This is not a soft delete! This action is irreversible.
	</p>
	<button
		class="text-left w-fit bg-red-400/30 border-[0.5px] border-red-400 px-2 py-1.5 text-red-400 hover:bg-red-400 hover:text-red-50 hover:cursor-pointer transition-colors duration-100 mt-4"
		onclick={() => (activateDelete = true)}
	>
		Delete account
	</button>
</div>

{#if activateDelete}
	<DestructiveBar buttonText="Delete Account" destructiveFn={handleDelete} bind:isSubmitting />
{/if}
