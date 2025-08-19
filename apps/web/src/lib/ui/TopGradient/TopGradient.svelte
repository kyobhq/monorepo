<script lang="ts">
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import { crossfade } from 'svelte/transition';

	const currentServer = $derived(serverStore.getServer(page.params.server_id || ''));

	let calculatedOpacity = $state(0.5); // Default opacity

	const [send, receive] = crossfade({
		duration: 300,
		fallback: (node) => {
			return {
				duration: 300,
				css: (t) => `opacity: ${t}`
			};
		}
	});

	/**
	 * Calculate the perceived lightness of an image
	 * @param imageUrl - The URL of the image to analyze
	 * @returns Promise that resolves to a lightness value between 0 (dark) and 1 (light)
	 */
	async function calculateImageLightness(imageUrl: string): Promise<number> {
		return new Promise((resolve, reject) => {
			const img = new Image();
			img.crossOrigin = 'anonymous';

			img.onload = () => {
				try {
					const canvas = document.createElement('canvas');
					const ctx = canvas.getContext('2d');

					if (!ctx) {
						reject(new Error('Could not get canvas context'));
						return;
					}

					// Use a smaller canvas for performance
					const size = Math.min(img.width, img.height, 100);
					canvas.width = size;
					canvas.height = size;

					// Draw and scale the image
					ctx.drawImage(img, 0, 0, size, size);

					// Get image data
					const imageData = ctx.getImageData(0, 0, size, size);
					const data = imageData.data; // looks like this: [r,g,b,a,r,g,b,a,...] == [pixel1, pixel2,...]

					let totalLightness = 0;
					let pixelCount = 0;

					// Sample every 4th pixel for performance
					for (let i = 0; i < data.length; i += 16) {
						const r = data[i];
						const g = data[i + 1];
						const b = data[i + 2];
						const a = data[i + 3];

						// Skip transparent pixels
						if (a < 128) continue;

						// Calculate perceived lightness using the relative luminance formula
						// This weighs green more heavily as the human eye is more sensitive to it
						const lightness = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
						totalLightness += lightness;
						pixelCount++;
					}

					if (pixelCount === 0) {
						resolve(0.5); // Default if no valid pixels
						return;
					}

					const averageLightness = totalLightness / pixelCount;
					resolve(averageLightness);
				} catch (error) {
					reject(error);
				}
			};

			img.onerror = () => reject(new Error('Failed to load image'));
			img.src = imageUrl;
		});
	}

	/**
	 * Map lightness value to optimal opacity
	 * @param lightness - Value between 0 (dark) and 1 (light)
	 * @returns Opacity value between 0.2 and 0.8
	 */
	function mapLightnessToOpacity(lightness: number): number {
		// Inverse relationship: darker images need more opacity for visibility
		// Light images need less opacity to avoid being too bright
		const minOpacity = 0.2;
		const maxOpacity = 1;

		// lerp that shit
		const opacity = maxOpacity - lightness * (maxOpacity - minOpacity);

		// bam
		return Math.max(minOpacity, Math.min(maxOpacity, opacity));
	}

	/**
	 * Update opacity based on current server avatar
	 */
	async function updateOpacity() {
		if (!currentServer?.avatar) {
			calculatedOpacity = 0.5;
			return;
		}

		try {
			const lightness = await calculateImageLightness(currentServer.avatar);
			calculatedOpacity = mapLightnessToOpacity(lightness);
		} catch (error) {
			console.warn('Failed to calculate image lightness:', error);
			calculatedOpacity = 0.5; // Fallback to default
		}
	}

	// Reactively update opacity when server changes
	$effect(() => {
		if (currentServer?.avatar) {
			updateOpacity();
		}
	});
</script>

{#if currentServer}
	<figure
		class="fixed -top-[7.5rem] left-0 w-screen h-[15rem] z-[10] pointer-events-none blur-3xl"
		style="opacity: {calculatedOpacity}"
	>
		{#key currentServer.avatar}
			<img
				class="absolute inset-0 fade-out-gradient bg-cover bg-no-repeat bg-center w-full h-full rounded-[50%] object-cover select-none mix-blend-plus-lighter"
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
