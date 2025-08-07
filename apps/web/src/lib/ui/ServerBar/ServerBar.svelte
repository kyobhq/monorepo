<script lang="ts">
	import Input from 'ui/Input/Input.svelte';
	import BarSeparator from 'ui/BarSeparator/BarSeparator.svelte';
	import MagnifyingGlass from 'ui/icons/MagnifyingGlass.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import type { Member } from '$lib/types/types';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
	import UserLine from 'ui/UserLine/UserLine.svelte';
	import { userStore } from 'stores/userStore.svelte';

	let membersPerRole = $derived.by(() => {
		if (!serverStore.roles || serverStore.roles.length === 0) {
			return {
				Members: serverStore.members
			};
		}

		let roleGroups: Record<string, Member[]> = {};

		for (const role of serverStore.roles) {
			roleGroups[role.name] = [];
		}
		roleGroups['Members'] = [];

		for (const member of serverStore.members) {
			let hasRole = false;

			for (const role of serverStore.roles) {
				if (member.roles.includes(role.id)) {
					roleGroups[role.name].push(member);
					hasRole = true;
					break;
				}
			}

			if (!hasRole) {
				roleGroups['Members'].push(member);
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

<div class="bg-main-975 border-l-[0.5px] border-l-main-700 w-[16rem] overflow-hidden">
	<section class="p-2.5">
		<Input Icon={MagnifyingGlass} placeholder="Search" />
	</section>
	<BarSeparator title={`Members - ${serverStore.memberCount}`} />
	<div class="flex flex-col gap-y-2 p-2.5">
		{#each Object.keys(membersPerRole) as role, idx (idx)}
			<CollapsibleBox header={role} canCollapse={false}>
				{#each membersPerRole[role] as member (member.id)}
					{@const isCurrentUser = member.id === userStore.user!.id}
					{@const displayName = isCurrentUser ? userStore.user!.display_name : member.display_name!}
					{@const avatar = isCurrentUser ? userStore.user!.avatar : member.avatar!}
					{@const status = isCurrentUser ? 'online' : member.status}

					<UserLine id={member.id!} {status} name={displayName} {avatar} hoverable />
				{/each}
			</CollapsibleBox>
		{/each}
	</div>
</div>
