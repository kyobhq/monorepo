<script lang="ts">
	import HashChat from 'ui/icons/HashChat.svelte';

	interface Props {
		name: string;
		description?: string;
	}

	let { name, description }: Props = $props();
</script>

<div class="absolute flex gap-x-1 py-2 px-2.5 w-fit z-10 select-none">
	<div
		class="flex border-[1.5px] border-main-800 bg-main-950/75 gap-x-2.5 px-3 py-2 items-center rounded-lg backdrop-blur-lg shadow-xs text-main-200"
	>
		<HashChat height={18} width={18} />
		{name}
	</div>
	{#if description}
		<div
			class="border-[1.5px] border-main-800 bg-main-950/75 gap-x-2.5 px-3 py-2 text-main-400 rounded-lg backdrop-blur-lg shadow-xs"
		>
			{description}
		</div>
	{/if}
</div>
