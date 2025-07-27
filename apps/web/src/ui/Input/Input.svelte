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
	class="flex items-center gap-3 bg-main-950 border pl-3 border-main-700 text-300 focus-within:border-main-500 transition duration-100"
>
	<Icon height={18} width={18} />
	<input {placeholder} {disabled} class="placeholder:text-main-500 w-full py-2 text-sm" />
</div>
