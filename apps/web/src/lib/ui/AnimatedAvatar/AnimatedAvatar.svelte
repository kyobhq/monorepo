<script lang="ts">
	import { onMount } from 'svelte';

	interface Props {
		src: string;
		alt?: string;
		class?: string;
		width?: number;
		height?: number;
		hover: boolean;
	}

	let { src, alt = '', class: className = '', width, height, hover }: Props = $props();

	let canvasEl = $state<HTMLCanvasElement>();
	let staticDataUrl = $state<string>('');
	let isLoaded = $state(false);

	onMount(() => {
		if (!canvasEl) return;

		const img = new Image();
		img.crossOrigin = 'anonymous';

		img.onload = () => {
			const canvas = canvasEl;
			if (!canvas) return;

			const ctx = canvas.getContext('2d');
			if (!ctx) return;

			canvas.width = img.naturalWidth;
			canvas.height = img.naturalHeight;

			ctx.drawImage(img, 0, 0);

			staticDataUrl = canvas.toDataURL('image/webp', 0.8);
			isLoaded = true;
		};

		img.onerror = () => {
			console.error('Failed to load animated avatar:', src);
			staticDataUrl = src;
			isLoaded = true;
		};

		img.src = src;
	});
</script>

<div
	class={className}
	role="img"
	aria-label={alt}
	style={width && height ? `width: ${width}px; height: ${height}px;` : ''}
>
	<canvas bind:this={canvasEl} class="hidden"></canvas>

	{#if isLoaded && staticDataUrl && !hover}
		<img
			src={staticDataUrl}
			{alt}
			class="w-full h-full object-cover"
			style="pointer-events: none;"
		/>
	{/if}

	{#if hover}
		<img {src} {alt} class="w-full h-full object-cover" style="pointer-events: none;" />
	{/if}

	{#if !isLoaded}
		<img {src} {alt} class="w-full h-full object-cover" style="pointer-events: none;" />
	{/if}
</div>
