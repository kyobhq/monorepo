<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { EditServerSchema } from '$lib/types/schemas';
	import FormInput from 'ui/Form/FormInput.svelte';
	import SaveBar from 'ui/SaveBar/SaveBar.svelte';
	import Switch from 'ui/Switch/Switch.svelte';
	import type { Server } from '$lib/types/types';
	import { backend } from 'stores/backendStore.svelte';
	import { logErr } from 'utils/print';
	import { serverStore } from 'stores/serverStore.svelte';
	import { generateTextWithExt } from 'utils/richInput';
	import { hasPermissions, isOwner } from 'utils/permissions';

	let { server }: { server: Server } = $props();
	let initialized = $state(false);
	let isSubmitting = $state(false);
	let isButtonComplete = $state(true);

	const { form, errors, enhance } = superForm(defaults(valibot(EditServerSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(EditServerSchema),
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isButtonComplete = false;
				isSubmitting = true;

				const res = await backend.editServerProfile(server.id, form.data);
				res.match(
					() => {
						serverStore.updateProfile(server.id, form.data);
					},
					(error) => logErr(error)
				);

				isSubmitting = false;
			}
		}
	});

	let changes = $derived.by(() => {
		return (
			$form.name !== server.name ||
			generateTextWithExt($form.description) !== generateTextWithExt(server.description) ||
			$form.public !== server.public ||
			!isButtonComplete
		);
	});

	$effect(() => {
		$form.name = server.name;
		$form.description = server.description;
		$form.public = server.public;
		initialized = true;
	});
</script>

<form method="post" use:enhance class="w-full flex flex-col gap-y-6 relative mt-6">
	<FormInput
		title="Name"
		id="name"
		type="text"
		placeholder="Name"
		bind:inputValue={$form.name}
		bind:error={$errors.name}
		class="w-full"
	/>

	<FormInput
		title="Description"
		id="description"
		type="rich"
		placeholder="Description"
		bind:inputValue={$form.description}
		bind:error={$errors.description}
		class="w-full"
	/>

	{#if isOwner(server.id)}
		<Switch
			active={$form.public}
			action={() => ($form.public = true)}
			reverse={() => ($form.public = false)}
			label="Public server"
		/>
	{/if}

	{#if initialized && changes}
		<SaveBar bind:isSubmitting bind:isButtonComplete />
	{/if}
</form>
