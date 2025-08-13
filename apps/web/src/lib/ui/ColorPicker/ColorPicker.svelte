<script lang="ts">
	import ColorPicker from 'svelte-awesome-color-picker';

	interface Props {
		color: string;
	}

	let { color = $bindable() }: Props = $props();
	let isOpen = $state(false);
</script>

<div class="relative flex gap-x-1">
	<button
		type="button"
		class="border-main-800 hocus:border-main-200 h-[2.5rem] w-[4rem] border transition-colors hover:cursor-pointer"
		aria-label="role color"
		style="background-color: {color}"
		onclick={() => (isOpen = !isOpen)}
	></button>
	{#if isOpen}
		<div class="dark absolute top-[3rem] -left-[12px] z-[50] [&_*]:!rounded-[6px]">
			<ColorPicker
				bind:hex={color}
				{isOpen}
				sliderDirection="vertical"
				isDialog={false}
				isDark={true}
				isAlpha={false}
				--picker-height="200px"
				--picker-width="200px"
				--picker-z-index="10"
			/>
		</div>
	{/if}
</div>

<style>
	.dark {
		--cp-bg-color: var(--color-main-800);
		--cp-border-color: var(--color-main-700);
		--cp-text-color: white;
		--cp-input-color: var(--color-main-900);
		--cp-button-hover-color: #777;
	}
</style>
