<script lang="ts">
	import DefaultDialog from '../DefaultDialog/DefaultDialog.svelte';
	import { AddFriendSchema } from '$lib/types/schemas';
	import { defaults, setError, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import FormInput from 'ui/Form/FormInput.svelte';
	import DialogFooter from '../DialogFooter/DialogFooter.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { backend } from 'stores/backendStore.svelte';
	import { logErr } from 'utils/print';
	import { userStore } from 'stores/userStore.svelte';

	const { form, errors, enhance } = superForm(defaults(valibot(AddFriendSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(AddFriendSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				const res = await backend.sendFriendRequest(form.data);
				res.match(
					(friend) => {
						coreStore.friendsDialog = false;
						userStore.friends.push(friend);
					},
					(err) => {
						setError(form, 'friend_username', err.message);
						setTimeout(() => {
							$errors = {};
						}, 4000);
						logErr(err);
					}
				);
			}
		}
	});
</script>

<DefaultDialog
	bind:state={coreStore.friendsDialog}
	title="Add a friend"
	subtitle="You just need his username!"
>
	<form method="post" use:enhance>
		<FormInput
			title="Friend username"
			id="friend-username"
			type="text"
			bind:error={$errors.friend_username}
			bind:inputValue={$form.friend_username}
			placeholder="johndoe123"
			class="mt-4 px-8"
		/>
		<DialogFooter buttonText="Add Friend" />
	</form>
</DefaultDialog>
