<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { EditPasswordSchema } from '$lib/types/schemas';
	import SaveBar from 'ui/SaveBar/SaveBar.svelte';
	import FormInput from 'ui/Form/FormInput.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { backend } from 'stores/backendStore.svelte';

	let initialized = $state(false);
	let isSubmitting = $state(false);
	let isButtonComplete = $state(true);
	let globalError = $state<string | null>(null);

	const { form, errors, enhance } = superForm(defaults(valibot(EditPasswordSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(EditPasswordSchema),
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isButtonComplete = false;
				isSubmitting = true;
				globalError = null;

				const res = await backend.updatePassword(form.data);
				res.match(
					() => {
						form.data.current = '';
						form.data.new = '';
						form.data.confirm = '';
					},
					(err) => {
						console.error(`${err.code}: ${err.message}`);
						globalError = err.message;

						setTimeout(() => {
							globalError = null;
						}, 5000);
					}
				);

				isSubmitting = false;
			}
		}
	});

	let changes = $derived.by(() => {
		return ($form.new !== '' && $form.confirm !== '' && $form.current !== '') || !isButtonComplete;
	});

	$effect(() => {
		if (userStore.user) {
			initialized = true;
		}
	});
</script>

<form method="post" use:enhance class="w-full relative">
	<p class={['font-medium select-none', globalError ? 'text-red-400' : '']}>
		Password
		{#if globalError}
			<span class="text-red-400">- {globalError}</span>
		{/if}
	</p>
	<div class="flex flex-col gap-y-6 mt-2.5">
		<FormInput
			title="Current Password"
			id="password"
			type="password"
			placeholder="Current Password"
			bind:inputValue={$form.current}
			bind:error={$errors.current}
			class="w-full"
		/>
		<FormInput
			title="New Password"
			id="new"
			type="password"
			placeholder="New Password"
			bind:inputValue={$form.new}
			bind:error={$errors.new}
			class="w-full"
		/>
		<FormInput
			title="Confirm Password"
			id="confirm"
			type="password"
			placeholder="Confirm Password"
			bind:inputValue={$form.confirm}
			bind:error={$errors.confirm}
			class="w-full"
		/>
	</div>

	{#if initialized && changes}
		<SaveBar bind:isSubmitting bind:isButtonComplete bind:isError={globalError} />
	{/if}
</form>
