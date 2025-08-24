<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import ModDialog from './ModDialog.svelte';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { KickUserSchema } from '$lib/types/schemas';
	import FormInput from 'ui/Form/FormInput.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { backend } from 'stores/backendStore.svelte';
	import { logErr } from 'utils/print';

	const member = $derived(
		serverStore.getMember(coreStore.modDialog.server_id, coreStore.modDialog.user_id)
	);
	const { form, errors, enhance } = superForm(defaults(valibot(KickUserSchema)), {
		id: 'kick-user-form',
		dataType: 'json',
		SPA: true,
		validators: valibot(KickUserSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				form.data.user_id = coreStore.modDialog.user_id;

				coreStore.modDialog.open = false;
				const res = await backend.kickUser(coreStore.modDialog.server_id, form.data);
				if (res.isErr()) logErr(res.error);
			}
		}
	});
</script>

<ModDialog
	title={`Kick ${member?.display_name}`}
	subtitle={`Are you sure you want to kick ${member?.display_name} from this server? They'll be able to join again with a valid invite.`}
	open={coreStore.modDialog.open && coreStore.modDialog.action === 'kick'}
	bind:openState={coreStore.modDialog.open}
>
	<form method="post" use:enhance>
		<FormInput
			title="Kick Reason"
			id="kick-reason"
			type="textarea"
			bind:error={$errors.reason}
			bind:inputValue={$form.reason}
			placeholder="Kick reason"
			class="mt-4"
			inputClass="w-full"
		/>

		<div class="flex justify-end mt-5">
			<button
				class="bg-red-400/20 border-[0.5px] border-red-400 px-2 py-1 hover:cursor-pointer hover:bg-red-400 transition-colors duration-75 text-red-400 hover:text-red-50 rounded-md"
			>
				Kick
			</button>
		</div>
	</form>
</ModDialog>
