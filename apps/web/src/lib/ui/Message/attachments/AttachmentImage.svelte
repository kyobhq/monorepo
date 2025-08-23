<script lang="ts">
	let { image, medias, hover } = $props();
	let isAnimated = $derived(image.url.includes('-animated.webp'));
</script>

<button
	class={[
		'pointer-events-auto relative select-none rounded-md overflow-hidden highlight-border',
		medias.length > 1 && 'aspect-square'
	]}
>
	<div role="img" aria-label="Attachment">
		{#if isAnimated}
			{#if hover}
				<img
					src={image.url}
					alt=""
					class="w-full h-full object-cover"
					style="pointer-events: none;"
				/>
			{:else}
				<img
					src={image.url.replace('-animated.webp', '.webp')}
					alt=""
					class="w-full h-full object-cover"
					style="pointer-events: none;"
				/>
			{/if}
		{:else}
			<img
				src={image.url}
				alt=""
				class="w-full h-full object-cover"
				style="pointer-events: none;"
			/>
		{/if}
	</div>
</button>
