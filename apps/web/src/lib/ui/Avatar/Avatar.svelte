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
	let style = $derived.by(() => {
		let str = '';

		if (width && height) {
			str += `width: ${width}px; height: ${height}px;`;
		}

		if (isAnimated) {
			str += ` background-image: url(${src}); background-size: cover; background-repeat: no-repeat;`;
		}

		return str;
	});
</script>

<div class={className} role="img" aria-label={alt} {style}>
	{#if !isAnimated || (isAnimated && !hover)}
		<img
			src={src.replaceAll('-animated.webp', '.webp')}
			{alt}
			class="w-full h-full object-cover"
			style="pointer-events: none;"
		/>
	{/if}
</div>
