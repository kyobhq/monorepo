<script lang="ts">
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import PermissionsMember from './PermissionsMember.svelte';
	import { backend } from 'stores/backendStore.svelte';
	import { logErr } from 'utils/print';
	import Switch from 'ui/Switch/Switch.svelte';
	import PermissionsRole from './PermissionsRole.svelte';

	let { form = $bindable() } = $props();

	const owner_id = $derived(serverStore.getOwnerID(page.params.server_id!));
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

	function switchPrivate() {
		form.private = true;
		form.users = [owner_id];
		form.roles = [];
	}

	function switchPublic() {
		form.private = false;
		form.users = [];
		form.roles = [];
	}
</script>

<div class="flex flex-col gap-y-4 h-full">
	<Switch
		active={form.private}
		action={switchPrivate}
		reverse={switchPublic}
		label="Private channel"
	/>
	<div
		class={[
			'flex flex-col gap-y-0.5 transition-opacity duration-100',
			!form.private ? 'opacity-50 pointer-events-none' : ''
		]}
	>
		<p class="text-sm text-main-400 font-medium">Roles</p>
		<ul
			class="flex flex-col gap-y-0.5 h-[12rem] bg-main-900 border border-main-800 overflow-auto p-2 rounded-2xl"
		>
			{#each roles as role (role.id)}
				<PermissionsRole name={role.name} color={role.color} id={role.id} bind:form />
			{:else}
				<li>No roles in this server.</li>
			{/each}
		</ul>
	</div>
	<div
		class={[
			'flex flex-col gap-y-2 transition-opacity duration-100',
			!form.private ? 'opacity-50 pointer-events-none' : ''
		]}
	>
		<p class="text-sm text-main-400 font-medium">Members</p>
		<ul
			class="flex flex-col gap-y-2 h-[12rem] bg-main-900 border border-main-800 overflow-auto p-2 rounded-2xl"
		>
			{#each members as member (member.id)}
				<PermissionsMember
					id={member.id}
					avatar={member.avatar}
					display_name={member.display_name}
					bind:form
				/>
			{/each}
			{#if serverStore.memberCount - 1 > members.length}
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
