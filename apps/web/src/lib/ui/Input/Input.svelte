<script lang="ts">
	import type { Component } from 'svelte';

	interface Props {
		disabled?: boolean;
		Icon?: Component;
		placeholder: string;
		class?: string;
	}

	const { placeholder, Icon, disabled = false, class: classes }: Props = $props();
</script>

<div
	class={[
		'box-style text-300 flex items-center gap-3 pl-3 transition duration-100 rounded-[10px] z-[1] relative focus-within:border-main-700',
		classes
	]}
>
	<Icon height={18} width={18} />
	<input
		{placeholder}
		{disabled}
		class="w-full border-none py-2 placeholder:text-main-500 bg-transparent text-sm px-0"
	/>
</div>
