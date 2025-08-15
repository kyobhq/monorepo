<script lang="ts">
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import { crossfade } from 'svelte/transition';

	const currentServer = $derived(serverStore.getServer(page.params.server_id || ''));

	const [send, receive] = crossfade({
		duration: 300,
		fallback: (node) => {
			return {
				duration: 300,
				css: (t) => `opacity: ${t}`
			};
		}
	});
</script>

{#if currentServer}
	<figure
		class="fixed -top-[7.5rem] left-0 w-screen h-[15rem] z-[10] pointer-events-none blur-3xl opacity-30"
	>
		{#key currentServer.avatar}
			<img
				class="absolute inset-0 fade-out-gradient bg-cover bg-no-repeat bg-center w-full h-full rounded-[50%] object-cover"
				src={currentServer.avatar}
				aria-hidden="true"
				alt=""
				in:receive={{ key: currentServer.avatar }}
				out:send={{ key: currentServer.avatar }}
			/>
		{/key}
	</figure>
{/if}

<style>
	.fade-out-gradient {
		mask-image: linear-gradient(
			180deg,
			#fafafa 0%,
			rgba(250, 250, 250, 0.35169) 48.08%,
			rgba(250, 250, 250, 0.144231) 79.33%,
			rgba(250, 250, 250, 0) 100%
		);
	}
</style>
