<script lang="ts">
	interface Props {
		height?: number;
		width?: number;
		class?: string;
	}

	let { class: classes, height = 24, width = 24 }: Props = $props();
</script>

<svg
	{width}
	{height}
	class={classes}
	viewBox="0 0 24 24"
	fill="none"
	xmlns="http://www.w3.org/2000/svg"
>
	<g fill="none" stroke="currentColor" stroke-linecap="round" stroke-width="1.5"
		><path stroke-linejoin="round" d="m9.5 7l-5 5l5 5" /><path
			d="M4.5 12h10c1.667 0 5 1 5 5"
			opacity=".5"
		/></g
	>
</svg>
