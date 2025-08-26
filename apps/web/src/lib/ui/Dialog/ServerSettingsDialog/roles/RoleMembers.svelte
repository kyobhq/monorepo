<script lang="ts">
	import { backend } from 'stores/backendStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import Check from 'ui/icons/Check.svelte';
	import { logErr } from 'utils/print';

	let { roleID } = $props();
	const members = $derived(serverStore.getMembers(coreStore.serverSettingsDialog.server_id));

	async function onChangeMember(ev: any, memberID: string) {
		const target = ev.target as HTMLInputElement;

		if (target.checked) {
			const res = await backend.addRoleMember(
				coreStore.serverSettingsDialog.server_id,
				roleID,
				memberID
			);
			res.match(
				() => {},
				(err) => logErr(err)
			);
		} else {
			const res = await backend.removeRoleMember(
				coreStore.serverSettingsDialog.server_id,
				roleID,
				memberID
			);
			res.match(
				() => {},
				(err) => logErr(err)
			);
		}
	}
</script>

{#snippet memberLine(
	id: string,
	avatar: string,
	display_name: string,
	username: string,
	hasRole: boolean
)}
	<li class="flex items-center justify-between">
		<div class="flex items-center gap-x-2">
			<figure class="h-8 w-8 rounded-lg overflow-hidden">
				<img src={avatar} alt="" />
			</figure>
			<p class="font-medium">
				{display_name}
			</p>
			<p class="text-sm text-main-500">
				{username}
			</p>
		</div>

		<input
			checked={hasRole}
			id="activate-member-{id}"
			type="checkbox"
			class="appearance-none h-0 w-0 opacity-0 peer"
			onchange={(ev) => onChangeMember(ev, id)}
		/>
		<label
			for="activate-member-{id}"
			class="h-5 w-5 bg-main-800 border-[0.5px] border-main-600 hover:cursor-pointer peer-checked:border-accent-lighter peer-checked:bg-accent transition-colors duration-75 flex justify-center items-center text-main-800 peer-checked:text-main-50 rounded-md"
		>
			<Check height={14} width={14} />
		</label>
	</li>
{/snippet}

<ul class="flex flex-col gap-y-4 mt-6">
	{#each members as member (member.id)}
		{@render memberLine(
			member.id!,
			member.avatar!,
			member.display_name!,
			member.username!,
			Boolean(member?.roles?.find((r) => r === roleID))
		)}
	{/each}
</ul>
