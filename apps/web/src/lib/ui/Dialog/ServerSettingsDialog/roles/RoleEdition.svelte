<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { CreateOrUpdateRoleSchema } from '$lib/types/schemas';
	import RoleDisplay from './RoleDisplay.svelte';
	import RoleNavbar from './RoleNavbar.svelte';
	import RolePermissions from './RolePermissions.svelte';
	import SaveBar from 'ui/SaveBar/SaveBar.svelte';
	import type { Role } from '$lib/types/types';
	import { backend } from 'stores/backendStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { logErr } from 'utils/print';
	import RoleMembers from './RoleMembers.svelte';

	interface Props {
		selectedRole: Role | null;
		initialized: boolean;
	}

	let { selectedRole = $bindable(), initialized = $bindable() }: Props = $props();
	const TABS = ['Display', 'Permissions', 'Members'];
	const TABS_DEFAULT = ['Permissions'];
	let activeTab = $state();
	let isSubmitting = $state(false);
	let isButtonComplete = $state(true);

	const { form, errors, enhance } = superForm(defaults(valibot(CreateOrUpdateRoleSchema)), {
		SPA: true,
		dataType: 'json',
		validators: valibot(CreateOrUpdateRoleSchema),
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isSubmitting = true;
				const serverID = coreStore.serverSettingsDialog.server_id;

				const res = await backend.createOrUpdateRole(serverID, form.data);
				res.match(
					(role) => {
						if (!selectedRole) return;
						selectedRole.id = role.id;
						selectedRole.name = role.name;
						selectedRole.color = role.color;
						selectedRole.abilities = role.abilities;
						selectedRole.position = role.position;
					},
					(err) => logErr(err)
				);

				isSubmitting = false;
			}
		}
	});

	let changes = $derived.by(() => {
		if (!initialized || !selectedRole) return false;

		return (
			$form.name !== selectedRole.name ||
			$form.color !== selectedRole.color ||
			!$form.abilities.every((ability) => selectedRole?.abilities.includes(ability)) ||
			!selectedRole?.abilities.every((ability) => $form.abilities.includes(ability)) ||
			!isButtonComplete
		);
	});

	$effect(() => {
		if (selectedRole) {
			$form.id = selectedRole.id;
			$form.name = selectedRole.name;
			$form.abilities = selectedRole.abilities;
			$form.position = selectedRole.position;
			$form.color = selectedRole.color;
			activeTab = selectedRole.name === 'Default Permissions' ? 'Permissions' : 'Display';
			initialized = true;
		}
	});
</script>

<div class="h-full w-full overflow-y-auto p-4">
	<RoleNavbar
		TABS={selectedRole?.name === 'Default Permissions' ? TABS_DEFAULT : TABS}
		bind:activeTab
	/>

	<form method="post" use:enhance>
		{#if activeTab === 'Display'}
			<RoleDisplay bind:form={$form} bind:selectedRole errors={$errors} />
		{:else if activeTab === 'Permissions'}
			<RolePermissions bind:form={$form} />
		{:else if activeTab === 'Members'}
			<RoleMembers roleID={$form.id} />
		{/if}

		{#if changes}
			<SaveBar bind:isSubmitting bind:isButtonComplete />
		{/if}
	</form>
</div>
