<script lang="ts">
	let { video } = $props();
</script>

<video controls class="pointer-events-auto h-full w-full object-cover">
	<source src={video.url} />
	<track kind="captions" />
</video>
