<script lang="ts">
	import type { Component } from 'svelte';

	interface Props {
		disabled?: boolean;
		Icon?: Component;
		placeholder: string;
	}

	const { placeholder, Icon, disabled = false }: Props = $props();
</script>

<div
	class="text-300 flex items-center gap-3 border-[0.5px] border-main-700 bg-main-950 pl-3 transition duration-100 focus-within:border-main-500"
>
	<Icon height={18} width={18} />
	<input
		{placeholder}
		{disabled}
		class="w-full border-none py-2 placeholder:text-main-500 bg-transparent text-sm px-0"
	/>
</div>
