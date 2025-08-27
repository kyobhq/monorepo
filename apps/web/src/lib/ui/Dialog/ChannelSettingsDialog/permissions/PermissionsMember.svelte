<script lang="ts">
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import Check from 'ui/icons/Check.svelte';

	let { id, avatar, display_name, form = $bindable() } = $props();
	const isOwner = $derived(serverStore.getOwnerID(page.params.server_id!) === id);
	const topRole = $derived(serverStore.getUserTopRole(page.params.server_id!, id));

	function handleChange(e: any) {
		const target = e.target as HTMLInputElement;

		if (target.checked) {
			form.users = [...form.users, id];
		} else {
			form.users = form.users.filter((userID) => userID !== id);
		}
	}
</script>

<li
	class={[
		'flex items-center justify-between hover:bg-main-800 transition-colors pl-1 pr-3 py-1 rounded-xl',
		isOwner ? 'opacity-50 pointer-events-none' : ''
	]}
>
	<div class="flex items-center gap-x-2">
		<figure class="h-8 w-8 rounded-lg overflow-hidden">
			<img src={avatar} alt="" />
		</figure>
		<p class="font-medium" style="color: {topRole ? topRole.color : '#fff'};">
			{display_name}
		</p>
	</div>

	{#if !isOwner}
		<input
			checked={form.users.includes(id)}
			id="activate-member-{id}"
			type="checkbox"
			class="appearance-none h-0 w-0 opacity-0 peer"
			onchange={handleChange}
		/>
		<label
			for="activate-member-{id}"
			class="h-5 w-5 bg-main-800 border-[0.5px] border-main-600 hover:cursor-pointer peer-checked:border-accent-lighter peer-checked:bg-accent transition-colors duration-75 flex justify-center items-center text-main-800 peer-checked:text-main-50 rounded-md"
		>
			<Check height={14} width={14} />
		</label>
	{:else}
		<p class="text-main-400">Owner</p>
	{/if}
</li>
