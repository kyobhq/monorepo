<script lang="ts">
	import { serverStore } from 'stores/serverStore.svelte';
	import { page } from '$app/state';
	import type { Member } from '$lib/types/types';
	import MemberLine from './MemberLine.svelte';
	import { backend } from 'stores/backendStore.svelte';
	import { logErr } from 'utils/print';
	import { onNavigate } from '$app/navigation';

	let idx = $state(0);
	let loadingMembers = $state(false);
	let canLoadMore = $state(
		serverStore.memberCount > serverStore.getMembers(page.params.server_id!).length
	);
	const MEMBERS_PER_PAGE = 50;
	const THRESHOLD = 100;

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

	async function fetchMembers(offset: number) {
		const res = await backend.getServerMembers(page.params.server_id || '', offset);
		res.match(
			(members) => {
				if (members.length < MEMBERS_PER_PAGE || !members?.length) canLoadMore = false;
				serverStore.addMembers(page.params.server_id || '', members);
			},
			(err) => logErr(err)
		);
	}

	async function handleScroll(e: Event) {
		const target = e.target as HTMLDivElement;
		const limit = target.scrollHeight - target.clientHeight - THRESHOLD;

		if (target.scrollTop > limit && !loadingMembers && canLoadMore) {
			loadingMembers = true;
			idx += 1;
			await fetchMembers(idx * MEMBERS_PER_PAGE);
			loadingMembers = false;
		}
	}

	onNavigate(({ from, to }) => {
		const fromServerID = from?.params?.server_id;
		const toServerID = to?.params?.server_id;

		if (fromServerID !== toServerID) {
			serverStore.resetMembersList(fromServerID!);
		}
	});
</script>

<div
	class="relative flex flex-col gap-y-2 px-2.5 pt-2.5 pb-[9rem] overflow-auto h-full"
	onscroll={handleScroll}
>
	{#if Object.keys(membersPerRole).length > 0}
		<div class="flex flex-col gap-y-1">
			{#each Object.keys(membersPerRole) as roleID, idx (idx)}
				{@const role = serverStore.getRole(page.params.server_id || '', roleID)!}
				<p
					class="text-sm text-main-500 w-fit px-2 py-0.5 rounded-md mt-2 first:mt-0"
					style="background-color: {role.color || '#a1a1aa'}26; color: {role.color || '#a1a1aa'}"
				>
					{role.name}
				</p>
				{#each membersPerRole[roleID] as member (member.id)}
					<MemberLine {member} {role} />
				{/each}
			{/each}
		</div>
	{/if}
</div>
