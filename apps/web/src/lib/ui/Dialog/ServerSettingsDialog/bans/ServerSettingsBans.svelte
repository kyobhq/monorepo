<script lang="ts">
	import type { Member } from '$lib/types/types';
	import { backend } from 'stores/backendStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { onMount } from 'svelte';
	import { logErr } from 'utils/print';

	let bannedMembers = $state<Member[]>();

	onMount(async () => {
		const res = await backend.getBannedMembers(coreStore.serverSettingsDialog.server_id);
		if (res.isErr()) logErr(res.error);

		if (res.isOk()) {
			bannedMembers = res.value;
		}
	});

	async function unbanMember(memberID: string) {
		const res = await backend.unbanUser(coreStore.serverSettingsDialog.server_id, memberID);
		res.match(
			() => {
				bannedMembers = bannedMembers?.filter((member) => member.id !== memberID);
			},
			(err) => logErr(err)
		);
	}
</script>

{#if bannedMembers}
	<div
		class="flex-1 overflow-y-auto border-[0.5px] border-main-800 bg-main-900 mt-6 p-3 rounded-md"
	>
		<ul class="flex flex-col gap-y-2">
			{#each bannedMembers as member (member.id)}
				<li class="flex items-center justify-between">
					<div class="flex items-center gap-x-3">
						<img src={member.avatar} class="w-11 h-11 object-cover rounded-sm" alt="" />
						<div class="flex flex-col gap-y-0.5">
							<p class="leading-none font-medium">{member.display_name}</p>
							<p class="text-sm text-main-400 leading-none">{member.username}</p>
						</div>
					</div>

					<button
						class="bg-accent/20 text-accent-lighter px-2 py-1 rounded-md border-[0.5px] border-accent hover:bg-accent hover:text-main-50 transition duration-100 active:scale-[0.95]"
						onclick={() => unbanMember(member.id!)}
					>
						Unban
					</button>
				</li>
			{/each}
		</ul>
	</div>
{/if}
