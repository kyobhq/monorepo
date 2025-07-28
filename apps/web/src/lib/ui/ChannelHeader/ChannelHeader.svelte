<script lang="ts">
	import HashChat from 'ui/icons/HashChat.svelte';

	interface Props {
		name: string;
		description?: string;
	}

	let { name, description }: Props = $props();
</script>

<div class="flex gap-x-1 py-2 px-2.5 w-fit bg-gradient-to-b from-main-975 to-transparent">
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
