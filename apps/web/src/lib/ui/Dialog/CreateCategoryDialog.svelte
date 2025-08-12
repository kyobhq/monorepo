<script lang="ts">
	import { CreateCategorySchema } from '$lib/types/schemas';
	import { backend } from 'stores/backendStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import FormInput from 'ui/Form/FormInput.svelte';
	import DefaultDialog from './DefaultDialog.svelte';
	import DialogFooter from './DialogFooter.svelte';
	import { page } from '$app/state';
	import { categoryStore } from 'stores/categoryStore.svelte';

	const { form, errors, enhance } = superForm(defaults(valibot(CreateCategorySchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(CreateCategorySchema),
		async onUpdate({ form }) {
			if (form.valid) {
				form.data.server_id = page.params.server_id || '';
				form.data.position = categoryStore.getLastPositionInServer(page.params.server_id || '');
				form.data.e2ee = false;

				const res = await backend.createCategory(form.data);
				res.match(
					() => {},
					(error) => {
						console.error(`${error.code}: ${error.message}`);
					}
				);

				coreStore.categoryDialog = false;
			}
		}
	});
</script>

<DefaultDialog
	bind:state={coreStore.categoryDialog}
	title="Create a category"
	subtitle="A nice way to organize your server!"
>
	<form method="post" use:enhance>
		<FormInput
			title="Category name"
			id="category-name"
			type="text"
			bind:error={$errors.name}
			bind:inputValue={$form.name}
			placeholder="Cool stuff"
			class="mt-4 px-8"
		/>

		<!-- <div class="px-8 mt-4"> -->
		<!-- 	<Switch -->
		<!-- 		active={$form.public} -->
		<!-- 		action={() => ($form.public = true)} -->
		<!-- 		reverse={() => ($form.public = false)} -->
		<!-- 		label="Make this server public (people will be able to see it)" -->
		<!-- 	/> -->
		<!-- </div> -->

		<DialogFooter buttonText="Create category" />
	</form>
</DefaultDialog>
