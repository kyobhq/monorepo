<script lang="ts">
	import PlusIcon from 'ui/icons/PlusIcon.svelte';

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
	class="h-[3.5625rem] w-[3.5625rem] flex justify-center items-center bg-main-900 hover:bg-main-800/80 border-[0.5px] border-main-700 aspect-square text-main-500 hover:text-main-200 hover:cursor-pointer transition-colors duration-75 rounded-r-[2px] rounded-l-md"
	for="file-attachement"
>
	<input
		id="file-attachement"
		type="file"
		class="absolute h-0 w-0 text-transparent opacity-0"
		onchange={onFile}
		multiple
	/>
	<PlusIcon height={20} width={20} />
</label>
