<script lang="ts">
	import Input from 'ui/Input/Input.svelte';
	import BarSeparator from 'ui/BarSeparator/BarSeparator.svelte';
	import MagnifyingGlass from 'ui/icons/MagnifyingGlass.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import type { Member } from '$lib/types/types';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
	import UserLine from 'ui/UserLine/UserLine.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { page } from '$app/state';

	let membersPerRole = $derived.by(() => {
		if (!page.params.server_id) return {};
		const members = serverStore.getMembers(page.params.server_id);
		const roles = serverStore
			.getRoles(page.params.server_id)
			.filter((role) => role.name !== 'Default Permissions');

		if (roles.length === 0) {
			return {
				default: members
			};
		}

		let roleGroups: Record<string, Member[]> = {};

		for (const role of roles) {
			roleGroups[role.id] = [];
		}
		roleGroups['default'] = [];

		for (const member of members) {
			let hasRole = false;

			for (const role of roles) {
				if (member.roles?.includes(role.id)) {
					roleGroups[role.id].push(member);
					hasRole = true;
					break;
				}
			}

			if (!hasRole) {
				roleGroups['default'].push(member);
			}
		}

		const filteredRoleGroups: Record<string, Member[]> = {};
		for (const [role, members] of Object.entries(roleGroups)) {
			if (members.length > 0) {
				filteredRoleGroups[role] = members;
			}
		}

		return filteredRoleGroups;
	});
</script>

<div class="bg-main-975 border-l-[0.5px] border-l-main-800 w-[16rem] overflow-hidden">
	<section class="p-2.5">
		<Input Icon={MagnifyingGlass} placeholder="Search" />
	</section>
	<BarSeparator title={`Members - ${serverStore.memberCount}`} />
	<div class="flex flex-col gap-y-2 p-2.5">
		{#each Object.keys(membersPerRole) as roleID, idx (idx)}
			{@const role = serverStore.getRole(page.params.server_id || '', roleID)!}
			<CollapsibleBox header={role.name} canCollapse={false} color={role.color}>
				{#each membersPerRole[roleID] as member (member.id)}
					{@const isCurrentUser = member.id === userStore.user!.id}
					{@const displayName = isCurrentUser ? userStore.user!.display_name : member.display_name!}
					{@const avatar = isCurrentUser ? userStore.user!.avatar : member.avatar!}
					{@const status = isCurrentUser ? 'online' : member.status}

					<UserLine
						id={member.id!}
						{status}
						name={displayName}
						{avatar}
						hoverable
						color={role.color}
					/>
				{/each}
			</CollapsibleBox>
		{/each}
	</div>
</div>
