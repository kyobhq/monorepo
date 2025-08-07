<script lang="ts">
	interface Props {
		height: number;
		width: number;
		class?: string | string[];
	}

	let { class: classes, height = 24, width = 24 }: Props = $props();
</script>

<svg
	{width}
	{height}
	class={classes}
	viewBox="0 0 10 10"
	fill="none"
	xmlns="http://www.w3.org/2000/svg"
>
	<g clip-path="url(#clip0_255_1965)">
		<path
			d="M5.41699 4.5835L9.16699 0.833496M9.16699 0.833496H6.94033M9.16699 0.833496V3.06016"
			stroke="currentColor"
			stroke-linecap="round"
			stroke-linejoin="round"
		/>
		<path
			opacity="0.5"
			d="M9.16683 5.00016C9.16683 6.96433 9.16683 7.94641 8.55641 8.55641C7.94683 9.16683 6.96433 9.16683 5.00016 9.16683C3.036 9.16683 2.05391 9.16683 1.4435 8.55641C0.833496 7.94683 0.833496 6.96433 0.833496 5.00016C0.833496 3.036 0.833496 2.05391 1.4435 1.4435C2.05433 0.833496 3.036 0.833496 5.00016 0.833496"
			stroke="currentColor"
			stroke-linecap="round"
		/>
	</g>
	<defs>
		<clipPath id="clip0_255_1965">
			<rect width="10" height="10" fill="currentColor" />
		</clipPath>
	</defs>
</svg>
