<script lang="ts">
	import SubmitButton from 'ui/SubmitButton/SubmitButton.svelte';
	import { flyBlur } from 'utils/transition';

	let { isSubmitting = $bindable(), isButtonComplete } = $props();
</script>

<div
	transition:flyBlur={{ duration: 150, y: 5 }}
	class="w-[calc(100%-5rem)] absolute -bottom-3 left-1/2 -translate-x-1/2 flex justify-between items-center bg-main-800 border-[0.5px] border-main-600 py-1.5 pr-1.5 pl-4"
>
	<p>It seems you have unsaved changes!</p>
	<SubmitButton text="Save" bind:isSubmitting bind:isComplete={isButtonComplete} />
</div>
