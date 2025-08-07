<script lang="ts">
	interface Props {
		height?: number;
		width?: number;
		class?: string;
	}

	let { height = 24, width = 24, class: classes }: Props = $props();
</script>

<svg {width} {height} viewBox="0 0 24 24" class={classes}>
	<path
		opacity="0.5"
		fill-rule="evenodd"
		clip-rule="evenodd"
		d="M8.10627 18.2468C5.29819 16.0833 2 13.5422 2 9.1371C2 4.27416 7.50016 0.825465 12 5.50063V20.5C11 20.5 10 19.7294 8.96173 18.9109C8.68471 18.6925 8.39814 18.4717 8.10627 18.2468Z"
		fill="currentColor"
	/>
	<path
		d="M15.0383 18.9109C17.9806 16.5914 22 14 22 9.1371C22 4.27416 16.4998 0.825465 12 5.50063V20.5C13 20.5 14 19.7294 15.0383 18.9109Z"
		fill="currentColor"
	/>
</svg>
