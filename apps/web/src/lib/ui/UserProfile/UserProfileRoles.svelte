<script lang="ts">
	import { page } from '$app/state';
	import type { User } from '$lib/types/types';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';

	let { user }: { user: User } = $props();
	const roleIDs = $derived.by(() => {
		if (user.id === userStore.user!.id) {
			return serverStore.getUserRoles(page.params.server_id!);
		}

		return serverStore.getMemberRoles(page.params.server_id!, user.id);
	});
</script>

{#if roleIDs}
	<div class="flex gap-1.5 flex-wrap mt-3">
		{#each roleIDs as roleID (roleID)}
			{@const role = serverStore.getRole(page.params.server_id!, roleID)}
			{#if role && role.name !== 'Default Permissions'}
				<div
					class="z-[4] text-sm px-2 py-[3px] rounded-md mix-blend-plus-lighter"
					style="background-color: {role.color}4d; color: {role.color}"
				>
					{role.name}
				</div>
			{/if}
		{/each}
	</div>
{/if}
