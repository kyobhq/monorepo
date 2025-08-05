<script lang="ts">
	import Input from 'ui/Input/Input.svelte';
	import BarSeparator from 'ui/BarSeparator/BarSeparator.svelte';
	import MagnifyingGlass from 'ui/icons/MagnifyingGlass.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import type { Member } from '$lib/types/types';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
	import UserLine from 'ui/UserLine/UserLine.svelte';

	let membersPerRole = $derived.by(() => {
		if (!serverStore.roles) {
			return {
				Online: serverStore.members
			};
		}

		let map: Record<string, Member[]> = {};

		for (const role of serverStore.roles) {
			map[role.name] = [];
		}

		for (const member of serverStore.members) {
			for (const role of serverStore.roles) {
				if (member.roles.includes(role.id)) {
					map[role.name].push(member);
					break;
				}
			}

			map['Online'].push(member);
		}

		return map;
	});
</script>

<div class="bg-main-975 border-l-[0.5px] border-l-main-700 w-[19.5rem] overflow-hidden">
	<section class="p-2.5">
		<Input Icon={MagnifyingGlass} placeholder="Search" />
	</section>
	<BarSeparator title={`Members - ${serverStore.memberCount}`} />
	<div class="flex flex-col gap-y-2 p-2.5">
		{#each Object.keys(membersPerRole) as role, idx (idx)}
			<CollapsibleBox header={role} canCollapse={false}>
				{#each membersPerRole[role] as member (member.id)}
					<UserLine name={member.display_name!} avatar={member.avatar!} hoverable />
				{/each}
			</CollapsibleBox>
		{/each}
	</div>
</div>
