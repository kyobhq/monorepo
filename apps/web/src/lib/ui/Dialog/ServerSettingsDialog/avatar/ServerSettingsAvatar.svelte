<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { EditAvatarSchema } from '$lib/types/schemas';
	import SaveBar from 'ui/SaveBar/SaveBar.svelte';
	import type { Server } from '$lib/types/types';
	import Cropper from 'svelte-easy-crop';
	import { backend } from 'stores/backendStore.svelte';
	import { logErr } from 'utils/print';
	import { serverStore } from 'stores/serverStore.svelte';

	let { server }: { server: Server } = $props();

	let initialized = $derived(false);
	let isSubmitting = $state(false);
	let isButtonComplete = $state(true);
	let image = $state<File | undefined>();
	let avatar = $state<string | undefined>();
	let banner = $state<string | undefined>();

	let cropBanner = $state({ x: 0, y: 0 });
	let cropAvatar = $state({ x: 0, y: 0 });
	let cropBannerPixels = $state({ x: 0, y: 0, height: 0, width: 0 });
	let cropAvatarPixels = $state({ x: 0, y: 0, height: 0, width: 0 });
	let zoomAvatar = $state(1);
	let zoomBanner = $state(1);
	let minZoomAvatar = $state(3);
	let minZoomBanner = $state(3);
	let maxZoomAvatar = $state(5);
	let maxZoomBanner = $state(5);

	const { form, errors, enhance } = superForm(defaults(valibot(EditAvatarSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(EditAvatarSchema),
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isButtonComplete = false;
				isSubmitting = true;

				const res = await backend.updateServerAvatarAndBanner(
					server.id,
					cropAvatarPixels,
					cropBannerPixels,
					form.data
				);
				res.match(
					(res) => {
						avatar = undefined;
						banner = undefined;
						cropAvatarPixels = { x: 0, y: 0, height: 0, width: 0 };
						cropBannerPixels = { x: 0, y: 0, height: 0, width: 0 };
						form.data.avatar = undefined;
						form.data.banner = undefined;

						serverStore.updateAvatar(server.id, res.avatar, res.banner);
					},
					(error) => logErr(error)
				);

				isSubmitting = false;
			}
		}
	});

	let changes = $derived.by(() => {
		return $form.avatar || $form.banner || !isButtonComplete;
	});

	$effect(() => {
		initialized = true;
	});

	function onFile(e: Event, type: 'banner' | 'avatar') {
		const target = e.target as HTMLInputElement;
		image = target.files?.[0];

		if (image) {
			if (type === 'banner') $form.banner = image;
			if (type === 'avatar') $form.avatar = image;

			const dataUrl = URL.createObjectURL(image);
			const img = new Image();

			img.onload = () => {
				const aspectImage = img.naturalWidth / img.naturalHeight;

				if (type === 'avatar') {
					const aspectAvatar = 1;
					minZoomAvatar = Math.max(aspectAvatar / aspectImage, aspectImage / aspectAvatar);
					zoomAvatar = minZoomAvatar;
				} else {
					const aspectBanner = 320 / 224;
					minZoomBanner = Math.max(aspectBanner / aspectImage, aspectImage / aspectBanner);
					zoomBanner = minZoomBanner;
				}

				URL.revokeObjectURL(dataUrl);
			};

			img.src = dataUrl;
			if (type === 'banner') banner = dataUrl;
			if (type === 'avatar') avatar = dataUrl;
		}
	}
</script>

<form method="post" use:enhance class="w-full relative mt-6" enctype="multipart/form-data">
	<p
		class={[
			'text-main-500 text-sm select-none',
			$errors.avatar || $errors.banner ? 'text-red-400' : ''
		]}
	>
		Avatar and banner
		{#if $errors.avatar}
			<span class="text-red-400">- {$errors.avatar}</span>
		{:else if $errors.banner}
			<span class="text-red-400">- {$errors.banner}</span>
		{/if}
		{#if avatar || banner}
			- <button
				onclick={() => {
					avatar = undefined;
					banner = undefined;
				}}
				class="text-accent hover:text-accent-lighter transition-colors duration-100 hover:cursor-pointer underline"
				>Reset</button
			>
		{/if}
	</p>

	<div class="flex mt-1.5 gap-x-2">
		<div class="relative h-[224px] w-[320px]">
			{#if banner}
				<Cropper
					image={banner}
					cropSize={{ height: 224, width: 320 }}
					cropShape="rect"
					showGrid={false}
					bind:crop={cropBanner}
					bind:zoom={zoomBanner}
					minZoom={minZoomBanner}
					maxZoom={maxZoomBanner}
					oncropcomplete={(e) => {
						cropBannerPixels = e.pixels;
					}}
				/>
			{:else}
				<input
					type="file"
					id="banner-profile"
					name="banner-profile"
					aria-label="user banner"
					onchange={(e) => onFile(e, 'banner')}
					class="absolute inset-0 text-transparent peer z-[4] hover:cursor-pointer"
				/>
				<figure class="highlight-border w-full h-full peer-hover:after:border-main-50/75">
					<img class="w-full h-full object-cover select-none" src={server.banner} alt="" />
				</figure>
			{/if}
		</div>

		<div class="relative h-[85px] w-[85px]">
			{#if avatar}
				<Cropper
					image={avatar}
					cropSize={{ height: 85, width: 85 }}
					cropShape="rect"
					showGrid={false}
					bind:crop={cropAvatar}
					bind:zoom={zoomAvatar}
					minZoom={minZoomAvatar}
					maxZoom={maxZoomAvatar}
					oncropcomplete={(e) => {
						cropAvatarPixels = e.pixels;
					}}
				/>
			{:else}
				<input
					type="file"
					id="avatar-profile"
					name="avatar-profile"
					aria-label="user avatar"
					onchange={(e) => onFile(e, 'avatar')}
					class="absolute inset-0 text-transparent peer z-[4] hover:cursor-pointer"
				/>
				<figure class="highlight-border w-full h-full peer-hover:after:border-main-50/75">
					<img class="w-full h-full object-cover select-none" src={server.avatar} alt="" />
				</figure>
			{/if}
		</div>
	</div>

	{#if initialized && changes}
		<SaveBar
			text="It seems your server's avatar/banner is not saved yet!"
			bind:isSubmitting
			bind:isButtonComplete
		/>
	{/if}
</form>
