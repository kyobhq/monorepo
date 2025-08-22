<script lang="ts">
	import { channelStore } from 'stores/channelStore.svelte';

	let { channelID } = $props();
</script>

{#if channelStore.messageCache[channelID]?.scrollHeight > 0}
	<div
		style="
						height: {channelStore.messageCache[channelID].scrollHeight}px; 
						flex-shrink: 0;"
		aria-hidden="true"
	></div>
{/if}
