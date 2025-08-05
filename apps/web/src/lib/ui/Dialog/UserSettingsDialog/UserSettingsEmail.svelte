<script lang="ts">
	import FormInput from 'ui/Form/FormInput.svelte';
	import { userStore } from 'stores/userStore.svelte';

	let email = $state('');
	let error = $state('');

	$effect(() => {
		if (userStore.user) {
			email = userStore.user.email;
		}
	});
</script>

<form method="post" class="w-full relative mt-7">
	<p class="font-medium select-none">Email</p>
	<FormInput
		title="Email"
		id="email"
		type="email"
		placeholder="Email"
		bind:inputValue={email}
		bind:error
		class="w-full mt-2.5"
	/>
</form>
