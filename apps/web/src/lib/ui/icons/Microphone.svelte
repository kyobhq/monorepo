<script lang="ts">
	import { scaleBlur } from 'utils/transition';

	interface Props {
		height: number;
		width: number;
		class?: string;
		mute?: boolean;
	}

	let { class: classes, height = 24, width = 24, mute = false }: Props = $props();
</script>

{#if mute}
	<svg
		xmlns="http://www.w3.org/2000/svg"
		{width}
		{height}
		class={['absolute -scale-x-100', classes]}
		viewBox="0 0 256 256"
		in:scaleBlur={{ duration: 100, delay: 50, blurAmount: 8, startScale: 0 }}
		out:scaleBlur={{ duration: 100, blurAmount: 8, startScale: 0 }}
	>
		<path
			fill="currentColor"
			d="M213.38 229.92a8 8 0 0 1-11.3-.54l-30.92-34A78.8 78.8 0 0 1 136 207.59V240a8 8 0 0 1-16 0v-32.4A80.11 80.11 0 0 1 48 128a8 8 0 0 1 16 0a64.07 64.07 0 0 0 64 64a63.4 63.4 0 0 0 32.21-8.68l-11.1-12.2A48 48 0 0 1 80 128V95.09L42.08 53.38a8 8 0 0 1 11.84-10.76l160 176a8 8 0 0 1-.54 11.3m-24.19-63.13a7.9 7.9 0 0 0 3.51.82a8 8 0 0 0 7.19-4.49A79.16 79.16 0 0 0 208 128a8 8 0 0 0-16 0a63.3 63.3 0 0 1-6.48 28.09a8 8 0 0 0 3.67 10.7m-27.33-29.22a8 8 0 0 0 13.88-4.57a50 50 0 0 0 .26-5V64a48 48 0 0 0-92-19.13a8 8 0 0 0 1.41 8.57Z"
		/>
	</svg>
{:else}
	<svg
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 256 256"
		{width}
		{height}
		class={['absolute', classes]}
		in:scaleBlur={{ duration: 100, delay: 50, blurAmount: 8, startScale: 0 }}
		out:scaleBlur={{ duration: 100, blurAmount: 8, startScale: 0 }}
	>
		<path
			fill="currentColor"
			d="M80 128V64a48 48 0 0 1 96 0v64a48 48 0 0 1-96 0m128 0a8 8 0 0 0-16 0a64 64 0 0 1-128 0a8 8 0 0 0-16 0a80.11 80.11 0 0 0 72 79.6V240a8 8 0 0 0 16 0v-32.4a80.11 80.11 0 0 0 72-79.6"
		/>
	</svg>
{/if}
