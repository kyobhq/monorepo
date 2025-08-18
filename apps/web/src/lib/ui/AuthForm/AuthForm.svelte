<script lang="ts">
	import IdentityIcon from 'ui/icons/IdentityIcon.svelte';
	import AuthFormInput from './AuthFormInput.svelte';
	import Envelop from 'ui/icons/Envelop.svelte';
	import type { SuperFormData, SuperFormErrors } from 'sveltekit-superforms/client';

	interface Props {
		type: 'signin' | 'signup';
		buttonText: string;
		form: SuperFormData<any>;
		errors: SuperFormErrors<any>;
		enhance: any;
	}

	let { buttonText, type, form, errors, enhance }: Props = $props();
</script>

<form class="flex flex-col gap-y-4 mt-4" use:enhance>
	{#if type === 'signup'}
		<div class="flex gap-x-2">
			<AuthFormInput
				id="username"
				label="Username"
				placeholder="johndoe123"
				IWidth={24}
				type="text"
				Icon={IdentityIcon}
				bind:formInput={$form.username}
				errorInput={$errors?.username?.[0]}
			/>
			<AuthFormInput
				id="displayName"
				label="Display Name"
				placeholder="John Doe"
				IWidth={24}
				type="text"
				Icon={IdentityIcon}
				bind:formInput={$form.display_name}
				errorInput={$errors?.display_name?.[0]}
			/>
		</div>
	{/if}
	<AuthFormInput
		id="email"
		label="Email"
		placeholder="johndoe@gmail.com"
		type="email"
		Icon={Envelop}
		IWidth={20}
		bind:formInput={$form.email}
		errorInput={$errors?.email?.[0]}
	/>
	<AuthFormInput
		id="password"
		label="Password"
		placeholder="⦁⦁⦁⦁⦁⦁⦁⦁⦁⦁⦁⦁"
		type="password"
		bind:formInput={$form.password}
		errorInput={$errors?.password?.[0]}
	/>

	<button
		class="bg-accent-darker py-2 text-lg font-semibold hover:cursor-pointer border border-accent hover:bg-accent transition-colors duration-75 rounded-lg"
	>
		{buttonText}
	</button>
</form>
