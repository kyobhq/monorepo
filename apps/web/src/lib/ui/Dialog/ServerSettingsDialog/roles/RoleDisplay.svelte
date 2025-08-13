<script lang="ts">
	import { backend } from 'stores/backendStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import DestructiveBar from 'ui/DestructiveBar/DestructiveBar.svelte';
	import FormInput from 'ui/Form/FormInput.svelte';
	import FormSwitch from 'ui/Form/FormSwitch.svelte';
	import { logErr } from 'utils/print';

	let { form = $bindable(), selectedRole = $bindable(), errors } = $props();
	let isSubmitting = $state(false);
	let activateDelete = $state(false);

	async function handleDelete() {
		isSubmitting = true;
		const res = await backend.deleteRole(coreStore.serverSettingsDialog.server_id, selectedRole.id);
		res.match(
			() => {
				serverStore.deleteRole(coreStore.serverSettingsDialog.server_id, selectedRole.id);
				selectedRole = null;
			},
			(err) => logErr(err)
		);
		isSubmitting = false;
	}
</script>

<div class="mt-4 flex flex-col">
	<FormInput
		id="role-name"
		title="Role Name"
		type="text"
		bind:inputValue={form.name}
		error={errors.name}
		placeholder="Administrator"
	/>

	<FormInput
		id="role-color"
		title="Role Color"
		type="color-picker"
		bind:inputValue={form.color}
		error={errors.color}
		placeholder="Administrator"
		class="mt-4"
	/>

	<FormSwitch
		active={false}
		action={() => {}}
		reverse={() => {}}
		label="Display role members separatly from online members"
		class="border-t-[0.5px] mt-5"
	/>

	<FormSwitch
		active={false}
		action={() => {}}
		reverse={() => {}}
		label="Allow anyone to @mention this role"
		description="Members with the permission &quot;Mention @everyone and All Roles&quot; will still be able to ping this role."
	/>

	<button
		type="button"
		class="bg-red-400/25 border-[0.5px] border-red-400 text-red-400 w-fit px-2.5 py-1 mt-6 hover:bg-red-400 hover:text-red-50 transition-colors duration-100"
		onclick={() => (activateDelete = true)}
	>
		Delete Role
	</button>

	{#if activateDelete}
		<DestructiveBar buttonText="Delete Role" destructiveFn={handleDelete} bind:isSubmitting />
	{/if}
</div>
