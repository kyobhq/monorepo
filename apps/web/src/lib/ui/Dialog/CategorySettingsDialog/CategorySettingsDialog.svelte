<script lang="ts">
	import { page } from '$app/state';
	import { EditCategorySchema } from '$lib/types/schemas';
	import { backend } from 'stores/backendStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import FormInput from 'ui/Form/FormInput.svelte';
	import SaveBar from 'ui/SaveBar/SaveBar.svelte';
	import SideBarSettings from 'ui/SideBar/SideBarSettings.svelte';
	import DefaultSettingsDialog from '../DefaultSettingsDialog/DefaultSettingsDialog.svelte';
	import { categoryStore } from 'stores/categoryStore.svelte';

	let currentCategory = $derived(
		categoryStore.getCategory(page.params.server_id!, coreStore.categorySettingsDialog.category_id)
	);

	let initialized = $state(false);
	let isSubmitting = $state(false);
	let isButtonComplete = $state(true);

	const { form, errors, enhance } = superForm(defaults(valibot(EditCategorySchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(EditCategorySchema),
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isButtonComplete = false;
				isSubmitting = true;

				if (!page.params.server_id || !currentCategory) return;

				form.data.server_id = page.params.server_id;

				const res = await backend.editCategory(
					coreStore.categorySettingsDialog.category_id,
					form.data
				);
				res.match(
					() => {
						currentCategory.name = form.data.name;
						currentCategory.users = form.data.users;
						currentCategory.roles = form.data.roles;
					},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);

				isSubmitting = false;
			}
		}
	});

	let changes = $derived.by(() => {
		if (!initialized) return false;

		return $form.name !== currentCategory?.name || !isButtonComplete;
	});

	$effect(() => {
		if (currentCategory) {
			$form.name = currentCategory.name;
			$form.users = currentCategory.users || [];
			$form.roles = currentCategory.roles || [];
		}
	});
</script>

<DefaultSettingsDialog bind:state={coreStore.categorySettingsDialog.open} bind:initialized>
	<SideBarSettings
		settings={['Overview', 'Permissions']}
		navigationFn={(setting) => (coreStore.categorySettingsDialog.section = setting)}
		activeSection={coreStore.categorySettingsDialog.section}
	/>
	<div class="flex flex-col w-full h-full px-8 py-6">
		<h3 class="text-xl font-semibold">{coreStore.categorySettingsDialog.section}</h3>
		<form method="post" use:enhance class="w-full h-full flex flex-col gap-y-6 relative mt-8">
			{#if coreStore.categorySettingsDialog.section === 'Overview'}
				<FormInput
					title="Category name"
					id="category-name"
					type="text"
					placeholder="Category name"
					bind:inputValue={$form.name}
					bind:error={$errors.name}
					class="w-full"
				/>
			{:else if coreStore.categorySettingsDialog.section === 'Permissions'}
				Permissions
			{/if}

			{#if changes}
				<SaveBar bind:isSubmitting bind:isButtonComplete />
			{/if}
		</form>
	</div>
</DefaultSettingsDialog>
