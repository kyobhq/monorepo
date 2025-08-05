<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import DefaultSettingsDialog from './DefaultSettingsDialog.svelte';
	import SideBarSettings from 'ui/SideBar/SideBarSettings.svelte';
	import FormInput from 'ui/Form/FormInput.svelte';
	import { page } from '$app/state';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { EditChannelSchema } from '$lib/types/schemas';
	import { flyBlur } from 'utils/transition';
	import { backend } from 'stores/backendStore.svelte';
	import SubmitButton from 'ui/SubmitButton/SubmitButton.svelte';
	import SaveBar from 'ui/SaveBar/SaveBar.svelte';

	let initialized = $state(false);
	let isSubmitting = $state(false);
	let isButtonComplete = $state(true);

	const { form, errors, enhance } = superForm(defaults(valibot(EditChannelSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(EditChannelSchema),
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isButtonComplete = false;
				isSubmitting = true;

				const res = await backend.editChannel(
					coreStore.channelSettingsDialog.channel_id,
					form.data
				);
				res.match(
					() => {},
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

		return !isButtonComplete;
	});
</script>

<DefaultSettingsDialog bind:state={coreStore.channelSettingsDialog.open} bind:initialized>
	<SideBarSettings
		settings={['My Account', 'Profile']}
		navigationFn={(setting) => (coreStore.channelSettingsDialog.section = setting)}
		activeSection={coreStore.channelSettingsDialog.section}
	/>
	<div class="flex flex-col w-full h-full px-8 py-6">
		<h3 class="text-xl font-semibold">{coreStore.channelSettingsDialog.section}</h3>
		<form method="post" use:enhance class="w-full h-full flex flex-col gap-y-6 relative mt-8">
			{#if coreStore.channelSettingsDialog.section === 'My Account'}
				My account
			{:else if coreStore.channelSettingsDialog.section === 'Profile'}
				Profile
			{/if}

			{#if changes}
				<SaveBar bind:isSubmitting {isButtonComplete} />
			{/if}
		</form>
	</div>
</DefaultSettingsDialog>
