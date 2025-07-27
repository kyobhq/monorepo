<script lang="ts">
	import type { Component } from 'svelte';

	interface Props {
		tabs: {
			Icon: Component;
			href: string;
		}[];
		activeTab?: string;
		onclick: (href: string) => void;
	}

	const { tabs, onclick, activeTab }: Props = $props();

	let buttonEl = $state<HTMLButtonElement | null>();
	let left = $derived.by(() => {
		if (!buttonEl) return;
		if (activeTab.includes(buttonEl.dataset.route)) return buttonEl.offsetLeft - 4;
	});
</script>

<div class="flex items-center p-1 border border-main-700 w-fit relative">
	{#each tabs as tab}
		{@const Icon = tab.Icon}

		<button
			bind:this={buttonEl}
			onclick={() => onclick(tab.href)}
			class={'p-[0.4rem] z-[1] group hover:cursor-pointer'}
			data-route={tab.href}
		>
			<Icon
				height={18}
				width={18}
				class={[
					'transition-colors duration-100',
					activeTab?.includes(tab.href) ? 'text-main-50' : 'text-main-300 group-hocus:text-main-100'
				]}
			/>
		</button>
	{/each}

	<div
		class="absolute top-1 left-1 h-[calc(100%-0.5rem)] aspect-square bg-main-900 transition"
		style="transform: translateX({left}px)"
	></div>
</div>
