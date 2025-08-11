<script lang="ts">
	let { image, medias } = $props();
</script>

<button
	class={[
		'pointer-events-auto relative select-none rounded-md overflow-hidden highlight-border',
		medias.length > 1 && 'aspect-square'
	]}
>
	<img src={image.url} alt="Attachment" class="h-full w-full object-cover" />
</button>
