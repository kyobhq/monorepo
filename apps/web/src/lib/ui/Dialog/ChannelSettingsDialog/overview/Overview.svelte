<script lang="ts">
	import FormInput from 'ui/Form/FormInput.svelte';

	let { form = $bindable(), errors = $bindable() } = $props();
</script>

<FormInput
	title="Channel name"
	id="channel-name"
	type="text"
	placeholder="Channel name"
	bind:inputValue={form.name}
	bind:error={errors.name}
	class="w-full"
/>

<FormInput
	title="Channel description"
	id="channel-description"
	type="text"
	placeholder="Channel description"
	bind:inputValue={form.description}
	bind:error={errors.description}
	class="w-full"
/>
