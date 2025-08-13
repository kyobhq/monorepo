<script lang="ts">
	import { page } from '$app/state';
	import { CreateChannelSchema } from '$lib/types/schemas';
	import { backend } from 'stores/backendStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import FormInput from 'ui/Form/FormInput.svelte';
	import Switch from 'ui/Switch/Switch.svelte';
	import DefaultDialog from '../DefaultDialog/DefaultDialog.svelte';
	import DialogFooter from '../DialogFooter/DialogFooter.svelte';

	const { form, errors, enhance } = superForm(defaults(valibot(CreateChannelSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(CreateChannelSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				form.data.type = form.data.e2ee ? 'textual-e2ee' : 'textual';
				form.data.server_id = page.params.server_id || '';
				form.data.category_id = coreStore.channelDialog.category_id;
				form.data.position = channelStore.getChannelsLastPositionInCategory(
					page.params.server_id || '',
					coreStore.channelDialog.category_id
				);
				coreStore.channelDialog.open = false;

				const res = await backend.createChannel(form.data);
				res.match(
					() => {},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);
			}
		}
	});

	$effect(() => {
		if (!$form.type || $form.type === '') {
			$form.type = 'textual';
		}
	});
</script>

<DefaultDialog
	bind:state={coreStore.channelDialog.open}
	title="Create a channel"
	subtitle="You don't really have a choice if you want to talk to people ;)"
>
	<form method="post" use:enhance>
		<FormInput
			title="Name"
			id="channel-name"
			type="text"
			bind:error={$errors.name}
			bind:inputValue={$form.name}
			placeholder="General"
			class="mt-4 px-8"
		/>

		<FormInput
			title="Description"
			id="description-name"
			type="text"
			bind:error={$errors.description}
			bind:inputValue={$form.description}
			placeholder="Here you can talk about everything!"
			class="mt-4 px-8"
		/>

		<div class="px-8 mt-10">
			<Switch
				active={$form.e2ee}
				action={() => ($form.e2ee = true)}
				reverse={() => ($form.e2ee = false)}
				label="End to end encryption*"
			/>
		</div>

		<p class="text-main-600 text-xs px-8 mt-4">
			*End to end encryption make all your messages in this channel unreadable for everyone outside
			this channel, even us. However be careful, since encryption happens on your side it makes some
			features unusable like: Search, Automod, Usage of bots, etc. We advise to use those channels
			for conversations that should stay fully private even if it means less features.
		</p>

		<DialogFooter buttonText="Create channel" />
	</form>
</DefaultDialog>
