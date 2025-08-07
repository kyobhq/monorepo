<script lang="ts">
	import { JoinServerSchema } from '$lib/types/schemas';
	import { backend } from 'stores/backendStore.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import FormInput from 'ui/Form/FormInput.svelte';
	import DialogFooter from '../DialogFooter.svelte';
	import { serverStore } from 'stores/serverStore.svelte';

	const { form, errors, enhance } = superForm(defaults(valibot(JoinServerSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(JoinServerSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				const inviteInput = form.data.invite_link?.trim() ?? '';
				const match = inviteInput.match(/(?:\/invite\/)?([A-Za-z0-9_-]{4,})$/);
				const inviteId = match ? match[1] : inviteInput;

				const res = await backend.joinServerWithInvite(inviteId);
				res.match(
					(server) => {
						serverStore.addServer(server);
					},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);
			}
		}
	});
</script>

<form method="post" use:enhance>
	<FormInput
		title="Invite link"
		id="invite-link"
		type="text"
		bind:error={$errors.invite_link}
		bind:inputValue={$form.invite_link}
		placeholder={`${import.meta.env.VITE_DOMAIN}/invite/...`}
		class="mt-4 px-8"
	/>
	<DialogFooter buttonText="Join server" />
</form>
