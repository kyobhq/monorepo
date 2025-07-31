<script lang="ts">
	let { buttonText } = $props();
</script>

<div class="border-t-[0.5px] border-t-main-700 flex justify-end p-3 mt-8">
	<button
		class="bg-accent-darker border-[0.5px] border-accent px-2 py-1 hover:cursor-pointer hocus:bg-accent transition-colors duration-75"
	>
		{buttonText}
	</button>
</div>
