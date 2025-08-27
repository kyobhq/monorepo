<script lang="ts">
	import { page } from '$app/state';
	import { EditChannelSchema } from '$lib/types/schemas';
	import { backend } from 'stores/backendStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import SaveBar from 'ui/SaveBar/SaveBar.svelte';
	import SideBarSettings from 'ui/SideBar/SideBarSettings.svelte';
	import DefaultSettingsDialog from '../DefaultSettingsDialog/DefaultSettingsDialog.svelte';
	import Overview from './overview/Overview.svelte';
	import Permissions from './permissions/Permissions.svelte';

	let currentChannel = $derived(
		channelStore.getChannel(page.params.server_id || '', coreStore.channelSettingsDialog.channel_id)
	);

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

	let changes = $derived.by(() => {
		if (!initialized) return false;

		return (
			$form.name !== currentChannel?.name ||
			$form.description !== currentChannel?.description ||
			!$form.users?.every((userID) => currentChannel.users?.includes(userID)) ||
			!$form.roles?.every((roleID) => currentChannel.roles?.includes(roleID)) ||
			$form.users?.length !== currentChannel.users?.length ||
			$form.roles?.length !== currentChannel.roles?.length ||
			!isButtonComplete
		);
	});

	$effect(() => {
		if (currentChannel) {
			if (!currentChannel.users) currentChannel.users = [];
			if (!currentChannel.roles) currentChannel.roles = [];

			$form.name = currentChannel.name;
			$form.description = currentChannel.description;
			$form.users = currentChannel.users || [];
			$form.roles = currentChannel.roles || [];
			$form.private = currentChannel.users.length > 0 || currentChannel.roles.length > 0;
		}
	});
</script>

<DefaultSettingsDialog bind:state={coreStore.channelSettingsDialog.open} bind:initialized>
	<SideBarSettings
		settings={['Overview', 'Permissions']}
		navigationFn={(setting) => (coreStore.channelSettingsDialog.section = setting)}
		activeSection={coreStore.channelSettingsDialog.section}
	/>
	<div class="flex flex-col w-full h-full px-8 py-6">
		<h3 class="text-xl font-semibold">{coreStore.channelSettingsDialog.section}</h3>
		<form method="post" use:enhance class="w-full h-full flex flex-col gap-y-6 relative pt-8">
			{#if coreStore.channelSettingsDialog.section === 'Overview'}
				<Overview bind:form={$form} bind:errors={$errors} />
			{:else if coreStore.channelSettingsDialog.section === 'Permissions'}
				<Permissions bind:form={$form} />
			{/if}

			{#if changes}
				<SaveBar bind:isSubmitting bind:isButtonComplete />
			{/if}
		</form>
	</div>
</DefaultSettingsDialog>
