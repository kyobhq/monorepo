<script lang="ts">
	import { messageStore } from 'stores/messageStore.svelte';

	let { onClickSave = $bindable() } = $props();
</script>

<span class="text-xs relative z-[1]">
	escape to <button
		onclick={() => messageStore.stopEditing()}
		class="text-accent underline hover:text-accent-lighter hover:cursor-pointer">cancel</button
	>, enter to
	<button
		class="text-accent underline hover:text-accent-lighter hover:cursor-pointer"
		onclick={onClickSave}>save</button
	>
</span>
