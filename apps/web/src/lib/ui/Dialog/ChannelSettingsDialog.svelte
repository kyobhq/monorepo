<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import DefaultSettingsDialog from './DefaultSettingsDialog.svelte';
	import SideBarSettings from 'ui/SideBar/SideBarSettings.svelte';
	import FormInput from 'ui/Form/FormInput.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { page } from '$app/state';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { EditChannelSchema } from '$lib/types/schemas';
	import { flyBlur } from 'utils/transition';
	import { backend } from 'stores/backendStore.svelte';
	import SubmitButton from 'ui/SubmitButton/SubmitButton.svelte';

	let currentChannel = $derived(
		channelStore.getChannel(page.params.server_id || '', coreStore.channelSettingsDialog.channel_id)
	);

	let isSubmitting = $state(false);
	let isButtonComplete = $state(false);

	const { form, errors, enhance } = superForm(defaults(valibot(EditChannelSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(EditChannelSchema),
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isSubmitting = true;

				if (!page.params.server_id || !currentChannel) return;

				form.data.server_id = page.params.server_id;

				const res = await backend.editChannel(
					coreStore.channelSettingsDialog.channel_id,
					form.data
				);
				res.match(
					() => {
						currentChannel.name = form.data.name;
						currentChannel.description = form.data.description || '';
						currentChannel.users = form.data.users;
						currentChannel.roles = form.data.roles;
					},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);

				isSubmitting = false;
			}
		}
	});

	let changes = $derived(
		$form.name !== currentChannel?.name ||
			$form.description !== currentChannel?.description ||
			!isButtonComplete
	);

	$effect(() => {
		if (currentChannel) {
			$form.name = currentChannel.name;
			$form.description = currentChannel.description;
			$form.users = currentChannel.users || [];
			$form.roles = currentChannel.roles || [];
		}
	});
</script>

<DefaultSettingsDialog bind:state={coreStore.channelSettingsDialog.open}>
	<SideBarSettings
		settings={['Overview', 'Permissions']}
		navigationFn={(setting) => (coreStore.channelSettingsDialog.section = setting)}
		activeSection={coreStore.channelSettingsDialog.section}
	/>
	<div class="flex flex-col w-full h-full px-8 py-6">
		<h3 class="text-xl font-semibold">{coreStore.channelSettingsDialog.section}</h3>
		<form method="post" use:enhance class="w-full h-full flex flex-col gap-y-6 relative mt-8">
			{#if coreStore.channelSettingsDialog.section === 'Overview'}
				<FormInput
					title="Channel name"
					id="channel-name"
					type="text"
					placeholder="Channel name"
					bind:inputValue={$form.name}
					bind:error={$errors.name}
					class="w-full"
				/>

				<FormInput
					title="Channel description"
					id="channel-description"
					type="text"
					placeholder="Channel description"
					bind:inputValue={$form.description}
					bind:error={$errors.description}
					class="w-full"
				/>
			{:else if coreStore.channelSettingsDialog.section === 'Permissions'}
				Permissions
			{/if}

			{#if changes}
				<div
					transition:flyBlur={{ duration: 150, y: 5 }}
					class="w-[calc(100%-20rem)] absolute -bottom-3 left-1/2 -translate-x-1/2 flex justify-between items-center bg-main-800 border-[0.5px] border-main-600 py-1.5 pr-1.5 pl-4"
				>
					<p>It seems you have unsaved changes!</p>
					<SubmitButton text="Save" bind:isSubmitting bind:isComplete={isButtonComplete} />
				</div>
			{/if}
		</form>
	</div>
</DefaultSettingsDialog>
