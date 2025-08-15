<script lang="ts">
	import FilledPlusIcon from 'ui/icons/FilledPlusIcon.svelte';

	let { attachments = $bindable() } = $props();

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
	class="h-full px-3 flex justify-center items-center text-main-500 hover:text-main-200 hover:cursor-pointer transition-colors duration-75 z-[1]"
	for="file-attachement"
>
	<input
		id="file-attachement"
		type="file"
		class="absolute h-0 w-0 text-transparent opacity-0"
		onchange={onFile}
		multiple
	/>
	<FilledPlusIcon />
</label>
