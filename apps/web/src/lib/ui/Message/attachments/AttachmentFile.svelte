<script lang="ts">
	import FileIcon from 'ui/icons/FileIcon.svelte';

	let { file } = $props();

	function download() {
		const anchor = document.createElement('a');
		anchor.href = file.url;
		anchor.download = file.file_name;
		document.body.appendChild(anchor);
		anchor.click();
		document.body.removeChild(anchor);
	}
</script>

<button
	class="group border-main-700 bg-main-800 hocus:bg-accent-100/15 hocus:border-accent-100 relative flex items-center gap-x-3 border py-2 pr-5 pl-3 text-left transition-colors select-none hover:cursor-pointer"
	onclick={download}
>
	<FileIcon
		height={36}
		width={36}
		class="text-accent-200 group-hocus:text-accent-100 transition-colors"
	/>
	<div class="flex flex-col">
		<p class="group-hocus:text-accent-50 transition-colors">{file.file_name}</p>
		<span class="text-main-300 group-hocus:text-accent-50/70 text-xs transition-colors">
			{file.file_size}
		</span>
	</div>
</button>
