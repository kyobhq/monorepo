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
	import gsap from 'gsap';
	import { coreStore } from 'stores/coreStore.svelte';

	const currentServer = $derived(serverStore.getServer(page.params.server_id!));
	let membersEl = $state<HTMLElement>();

	let membersPerRole = $derived.by(() => {
		if (!page.params.server_id) return {};
		const members = serverStore.getMembers(page.params.server_id);
		const roles = serverStore
			.getRoles(page.params.server_id)
			.filter((role) => role.name !== 'Default Permissions');

		let roleGroups: Record<string, Member[]> = {};

		for (const role of roles) {
			roleGroups[role.id] = [];
		}
		roleGroups['default'] = [];
		roleGroups['offline'] = [];

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
				if (member.status === 'offline') roleGroups['offline'].push(member);
				else roleGroups['default'].push(member);
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

	$effect(() => {
		if (!membersEl || coreStore.firstLoad.serverbar) return;
		const boxes = membersEl.querySelectorAll('.collapsible-wrapper');
		if (!boxes.length) return;
		gsap.from(boxes, {
			opacity: 0,
			y: 8,
			stagger: 0.06,
			duration: 0.35,
			ease: 'power2.out',
			onComplete: () => {
				coreStore.firstLoad.serverbar = true;
			}
		});
	});
</script>

{#if currentServer && Object.keys(membersPerRole).length > 0}
	<div class="bg-main-975 border-l-[0.5px] border-l-main-800 w-[16rem] overflow-hidden h-screen">
		<section class="p-2.5">
			<Input Icon={MagnifyingGlass} placeholder={`Search ${currentServer.name}`} />
		</section>
		<BarSeparator title={`Members - ${serverStore.memberCount}`} />
		<div
			class="flex flex-col gap-y-2 p-2.5 h-[calc(100%-6.8rem)] overflow-auto"
			bind:this={membersEl}
		>
			{#each Object.keys(membersPerRole) as roleID, idx (idx)}
				{@const role = serverStore.getRole(page.params.server_id || '', roleID)!}
				<div class="collapsible-wrapper">
					<CollapsibleBox header={role.name} canCollapse={false} color={role.color}>
						{#each membersPerRole[roleID] as member (member.id)}
							{@const isCurrentUser = member.id === userStore.user!.id}
							{@const displayName = isCurrentUser
								? userStore.user!.display_name
								: member.display_name!}
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
				</div>
			{/each}
		</div>
	</div>
{/if}
