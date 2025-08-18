<script lang="ts">
	import { goto } from '$app/navigation';
	import { SignUpSchema } from '$lib/types/schemas';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import AuthForm from 'ui/AuthForm/AuthForm.svelte';
	import DangerIcon from 'ui/icons/DangerIcon.svelte';

	let globalError = $state<string>('');

	const { form, errors, enhance } = superForm(defaults(valibot(SignUpSchema)), {
		SPA: true,
		validators: valibot(SignUpSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				try {
					const res = await fetch(`${import.meta.env.VITE_API_URL}/signup`, {
						method: 'post',
						credentials: 'include',
						headers: {
							'Content-Type': 'application/json'
						},
						body: JSON.stringify({
							email: form.data.email,
							username: form.data.username,
							display_name: form.data.display_name,
							password: form.data.password
						})
					});

					if (!res.ok) {
						const data = await res.json();
						console.error('signup failed', res.status, data);
						throw new Error(`signup failed`);
					}

					return goto('/');
				} catch (err) {
					console.error(err);
					globalError = 'Signup failed';
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
	<div class="w-[26rem] absolute top-[30%] left-[12rem]">
		<h1 class="text-5xl font-bold text-main-50">Create an account.</h1>
		<p class="mt-4 text-main-400">
			Already have an account? <a href="/signin" class="text-accent-lighter">Sign in</a>
		</p>

		<AuthForm buttonText="Connect" type="signup" {enhance} {form} {errors} />
		{#if globalError !== ''}
			<div class="mt-6 flex items-center justify-center text-red-400 gap-x-3 -ml-6">
				<DangerIcon />
				<p class="font-medium">{globalError}</p>
			</div>
		{/if}
	</div>
</main>
