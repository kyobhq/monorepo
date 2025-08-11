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
		class={[
			'flex border-[0.5px] border-main-800 bg-main-975/80 gap-x-2.5 px-3 py-2 items-center rounded-l-md backdrop-blur-lg',
			description ? 'rounded-r-[2px]' : 'rounded-r-md'
		]}
	>
		<HashChat height={18} width={18} />
		{name}
	</div>
	{#if description}
		<div
			class="border-[0.5px] border-main-800 bg-main-975/80 gap-x-2.5 px-3 py-2 text-main-600 rounded-l-[2px] rounded-r-md backdrop-blur-lg"
		>
			{description}
		</div>
	{/if}
</div>
