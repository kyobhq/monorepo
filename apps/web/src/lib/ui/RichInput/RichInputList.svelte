<script lang="ts">
	import { editorStore } from 'stores/editorStore.svelte';
	import MentionsList from './extensions/mentions/MentionsList.svelte';
	import EmojisList from './extensions/emojis/EmojisList.svelte';
</script>

{#if editorStore.mentionProps}
	<MentionsList
		props={editorStore.mentionProps}
		bind:this={editorStore.mentionsListEl}
		class="w-full"
	/>
{:else if editorStore.emojiProps}
	<EmojisList props={editorStore.emojiProps} bind:this={editorStore.emojisListEl} class="w-full" />
{/if}
