<script lang="ts">
	import HashChat from 'ui/icons/HashChat.svelte';

	interface Props {
		name: string;
		description?: string;
	}

	let { name, description }: Props = $props();
</script>

<div class="absolute flex gap-x-1 py-2 px-2.5 w-fit z-10 backdrop select-none">
	<div class="flex border-[0.5px] border-main-800 bg-main-975 gap-x-2.5 px-3 py-2 items-center">
		<HashChat height={18} width={18} />
		{name}
	</div>
	{#if description}
		<div class="border-[0.5px] border-main-800 bg-main-975 gap-x-2.5 px-3 py-2 text-main-600">
			{description}
		</div>
	{/if}
</div>

<style>
	.backdrop {
		filter: drop-shadow(8px 0px 4px #121214) drop-shadow(8px 6px 4px #121214)
			drop-shadow(-8px -12px 4px #121214) drop-shadow(0px 6px 4px #121214);
	}
</style>
