<script lang="ts">
	interface Props {
		src: string;
		alt?: string;
		class?: string;
		width?: number;
		height?: number;
		hover: boolean;
	}

	let { src, alt = '', class: className = '', width, height, hover }: Props = $props();

	let isAnimated = $derived(src.includes('-animated.webp'));
</script>

<div
	class={className}
	role="img"
	aria-label={alt}
	style={width && height ? `width: ${width}px; height: ${height}px;` : ''}
>
	{#if isAnimated}
		{#if hover}
			<img {src} {alt} class="w-full h-full object-cover" style="pointer-events: none;" />
		{:else}
			<img
				src={src.replace('-animated.webp', '.webp')}
				{alt}
				class="w-full h-full object-cover"
				style="pointer-events: none;"
			/>
		{/if}
	{:else}
		<img {src} {alt} class="w-full h-full object-cover" style="pointer-events: none;" />
	{/if}
</div>
