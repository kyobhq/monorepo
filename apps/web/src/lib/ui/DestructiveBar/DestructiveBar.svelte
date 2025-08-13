<script lang="ts">
	import SubmitButton from 'ui/SubmitButton/SubmitButton.svelte';
	import { flyBlur } from 'utils/transition';

	let {
		buttonText = 'Delete',
		isSubmitting = $bindable(),
		isButtonComplete = $bindable(),
		isError = $bindable(),
		destructiveFn
	} = $props();
</script>

<div
	transition:flyBlur={{ duration: 150, y: 5 }}
	class="w-[calc(100%-5rem)] fixed -bottom-6 left-1/2 -translate-x-1/2 flex justify-between items-center bg-main-800 border-[0.5px] border-main-600 py-1.5 pr-1.5 pl-4 z-[1] rounded-md"
>
	<p class="select-none">Are you sure you want to continue with this action?</p>
	<SubmitButton
		type="button"
		text={buttonText}
		fn={destructiveFn}
		bind:isSubmitting
		bind:isComplete={isButtonComplete}
		bind:isError
		destructive
	/>
</div>
