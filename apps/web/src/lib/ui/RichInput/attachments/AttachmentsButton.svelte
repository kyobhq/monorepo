<script lang="ts">
	import FilledPlusIcon from 'ui/icons/FilledPlusIcon.svelte';

	let { attachments = $bindable(), disabled = false } = $props();

	function onFile(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files) {
			for (const file of target.files) {
				if (file.size > 10 << 20) {
					// errorsStore.attachmentError = true;
					break;
				}
				attachments.push(file);
			}
		}
	}
</script>

<label
	class={[
		'h-[3.5rem] px-3 flex justify-center items-center transition-colors duration-75 z-[1]',
		disabled
			? 'text-main-700 hover:cursor-not-allowed'
			: 'text-main-500 hover:text-main-200 hover:cursor-pointer'
	]}
	for="file-attachement"
>
	<input
		id="file-attachement"
		type="file"
		class="absolute h-0 w-0 text-transparent opacity-0"
		onchange={onFile}
		multiple
		{disabled}
	/>
	<FilledPlusIcon />
</label>
