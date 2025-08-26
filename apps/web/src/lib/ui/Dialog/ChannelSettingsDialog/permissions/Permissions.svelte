<script lang="ts">
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import PermissionsMember from './PermissionsMember.svelte';
	import { backend } from 'stores/backendStore.svelte';
	import { logErr } from 'utils/print';

	const roles = $derived(serverStore.getRoles(page.params.server_id!, true));
	const members = $derived(serverStore.getMembers(page.params.server_id!));

	async function loadMoreMembers() {
		const offset = serverStore.memberCount - members.length;

		const res = await backend.getServerMembers(page.params.server_id!, offset);
		res.match(
			(m) => {
				members.push(...m);
			},
			(err) => logErr(err)
		);
	}
</script>

<div class="flex flex-col gap-y-4 h-full">
	<div class="flex flex-col gap-y-0.5">
		<p class="text-sm text-main-400 font-medium">Roles</p>
		<ul>
			{#each roles as role (role.id)}
				<li>{role.name}</li>
			{:else}
				<li>No roles in this server.</li>
			{/each}
		</ul>
	</div>
	<div class="flex flex-col gap-y-2">
		<p class="text-sm text-main-400 font-medium">Members</p>
		<ul
			class="flex flex-col gap-y-2 h-[20rem] bg-main-900 border border-main-800 overflow-auto p-2 rounded-2xl"
		>
			{#each members as member (member.id)}
				<PermissionsMember
					id={member.id}
					avatar={member.avatar}
					display_name={member.display_name}
				/>
			{/each}
			{#if serverStore.memberCount > members.length}
				<button
					class="bg-main-800 border border-main-700 py-2 rounded-xl hover:bg-main-700/70 transition-colors duration-100 active-scale-down"
					type="button"
					onclick={loadMoreMembers}
				>
					Load more
				</button>
			{/if}
		</ul>
	</div>
</div>
