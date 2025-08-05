<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { EditUserSchema, type EditUserType } from '$lib/types/schemas';
	import SaveBar from 'ui/SaveBar/SaveBar.svelte';
	import FormInput from 'ui/Form/FormInput.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { backend } from 'stores/backendStore.svelte';
	import UserSettingsLinks from './UserSettingsLinks.svelte';
	import UserSettingsFacts from './UserSettingsFacts.svelte';
	import { generateTextWithExt } from 'utils/richInput';

	let initialized = $state(false);
	let isSubmitting = $state(false);
	let isButtonComplete = $state(true);

	const { form, errors, enhance } = superForm(defaults(valibot(EditUserSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(EditUserSchema),
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isButtonComplete = false;
				isSubmitting = true;

				const payload: EditUserType = {
					username: form.data.username,
					display_name: form.data.display_name,
					about_me: form.data.about_me,
					links: form.data.links,
					facts: form.data.facts
				};

				const res = await backend.updateProfile(payload);
				res.match(
					() => {
						if (!userStore.user) return;

						userStore.user.username = form.data.username;
						userStore.user.display_name = form.data.display_name;
						userStore.user.about_me = form.data.about_me;
						userStore.user.links = form.data.links;
						userStore.user.facts = form.data.facts;
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
		const userLinks = userStore.user?.links || [];
		const formLinks = $form.links || [];
		const userFacts = userStore.user?.facts || [];
		const formFacts = $form.facts || [];

		const linksChanged =
			userLinks.length !== formLinks.length ||
			userLinks.some((userLink, index) => {
				const formLink = formLinks[index];
				return !formLink || userLink.label !== formLink.label || userLink.url !== formLink.url;
			});

		const factsChanged =
			userFacts.length !== formFacts.length ||
			userFacts.some((userFact, index) => {
				const formFact = formFacts[index];
				return !formFact || userFact.label !== formFact.label || userFact.value !== formFact.value;
			});

		return (
			$form.username !== userStore.user?.username ||
			$form.display_name !== userStore.user?.display_name ||
			generateTextWithExt($form.about_me) !== generateTextWithExt(userStore.user?.about_me) ||
			linksChanged ||
			factsChanged ||
			!isButtonComplete
		);
	});

	$effect(() => {
		if (!initialized) {
			// Initialize with empty arrays if no user data yet
			$form.links = [];
			$form.facts = [];
		}

		if (userStore.user) {
			$form.username = userStore.user.username;
			$form.display_name = userStore.user.display_name;
			$form.about_me = userStore.user.about_me;
			$form.links = userStore.user.links || [];
			$form.facts = userStore.user.facts || [];
			initialized = true;
		}
	});
</script>

<form method="post" use:enhance class="w-full flex flex-col gap-y-6 relative mt-6">
	<FormInput
		title="Username"
		id="username"
		type="text"
		placeholder="Username"
		bind:inputValue={$form.username}
		bind:error={$errors.username}
		class="w-full"
	/>

	<FormInput
		title="Display Name"
		id="display-name"
		type="text"
		placeholder="Display Name"
		bind:inputValue={$form.display_name}
		bind:error={$errors.display_name}
		class="w-full"
	/>

	<FormInput
		title="About me"
		id="about-me"
		type="rich"
		placeholder="About me"
		bind:inputValue={$form.about_me}
		bind:error={$errors.about_me}
		class="w-full"
	/>

	<UserSettingsLinks bind:links={$form.links} />

	<UserSettingsFacts bind:facts={$form.facts} />

	{#if initialized && changes}
		<SaveBar bind:isSubmitting bind:isButtonComplete />
	{/if}
</form>
