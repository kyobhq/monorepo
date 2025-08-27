<script lang="ts">
	import { backend } from 'stores/backendStore.svelte';
	import Check from 'ui/icons/Check.svelte';
	import { logErr } from 'utils/print';

	let { id, name, color, form = $bindable() } = $props();

	async function handleChange(e: any) {
		const target = e.target as HTMLInputElement;
		let roleMemberIDs: string[] = [];

		const res = await backend.getRoleMembers(id);
		res.match(
			(memberIDs) => {
				roleMemberIDs = memberIDs;
			},
			(err) => logErr(err)
		);

		if (target.checked) {
			form.roles = [...form.roles, id];
			form.users = [...form.users, ...roleMemberIDs];
		} else {
			form.roles = form.roles.filter((roleID: string) => roleID !== id);
			form.users = form.users.filter((userID: string) => !roleMemberIDs.includes(userID));
		}
	}
</script>

<li class="flex items-center justify-between pl-1 pr-3 py-1 rounded-xl">
	<div
		class="flex items-center gap-x-2 py-1 px-2 rounded-md text-sm"
		style="color: {color}; background-color: {color}26;"
	>
		<p class="font-medium">{name}</p>
	</div>

	<input
		checked={form.roles.includes(id)}
		id="activate-role-{id}"
		type="checkbox"
		class="appearance-none h-0 w-0 opacity-0 peer"
		onchange={handleChange}
	/>
	<label
		for="activate-role-{id}"
		class="h-5 w-5 bg-main-800 border-[0.5px] border-main-600 hover:cursor-pointer peer-checked:border-accent-lighter peer-checked:bg-accent transition-colors duration-75 flex justify-center items-center text-main-800 peer-checked:text-main-50 rounded-md"
	>
		<Check height={14} width={14} />
	</label>
</li>
