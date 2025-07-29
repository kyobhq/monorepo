<script lang="ts">
	import { goto } from '$app/navigation';
	import { SignInSchema } from '$lib/types/schemas';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import AuthForm from 'ui/AuthForm/AuthForm.svelte';
	import DangerIcon from 'ui/icons/DangerIcon.svelte';

	let globalError = $state('');

	const { form, errors, enhance } = superForm(defaults(valibot(SignInSchema)), {
		SPA: true,
		validators: valibot(SignInSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				try {
					const res = await fetch(`${import.meta.env.VITE_API_URL}/signin`, {
						method: 'post',
						credentials: 'include',
						headers: {
							'Content-Type': 'application/json'
						},
						body: JSON.stringify({
							email: form.data.email,
							password: form.data.password
						})
					});

					if (!res.ok) {
						const data = await res.json();
						console.error('signin failed', res.status, data);
						globalError = data.error;
						return;
					}

					return goto('/');
				} catch (err) {
					console.error(err);
					globalError = 'Signin failed';
				}
			}
		}
	});
</script>

<img
	src="/images/bg.jpg"
	class="fixed top-0 left-0 h-screen w-screen object-cover -scale-x-100"
	alt=""
/>

<main class="relative h-screen w-screen auth-gradient z-[10]">
	<div class=" min-w-[26rem] absolute top-[30%] left-[12rem]">
		<h1 class="text-5xl font-bold text-main-50">Sign In.</h1>
		<p class="mt-4 text-main-400">
			No account yet? <a href="/signup" class="text-accent">Create one</a>
		</p>

		<AuthForm buttonText="Connect" type="signin" {enhance} {errors} {form} />
		{#if globalError !== ''}
			<div class="mt-6 flex items-center justify-center text-red-400 gap-x-3 -ml-6">
				<DangerIcon />
				<p class="font-medium">{globalError}</p>
			</div>
		{/if}
	</div>
</main>
