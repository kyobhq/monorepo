<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { AddEmojisSchema } from '$lib/types/schemas';
	import type { Emoji } from '$lib/types/types';
	import { transformShortcode } from 'utils/emojis';
	import { createId } from '@paralleldrive/cuid2';
	import { fly } from 'svelte/transition';
	import { userStore } from 'stores/userStore.svelte';
	import SubmitButton from 'ui/SubmitButton/SubmitButton.svelte';
	import DangerIcon from 'ui/icons/DangerIcon.svelte';
	import EmojiLine from 'ui/EmojiLine/EmojiLine.svelte';
	import { backend } from 'stores/backendStore.svelte';
	import { logErr } from 'utils/print';

	let emojis = $state<Emoji[]>([]);
	let isDeleting = $state();
	let isSubmitting = $state(false);
	let isButtonComplete = $state(true);

	const { form, errors, enhance } = superForm(defaults(valibot(AddEmojisSchema)), {
		SPA: true,
		dataType: 'json',
		validators: valibot(AddEmojisSchema),
		validationMethod: 'onsubmit',
		async onUpdate({ form }) {
			if (form.valid) {
				isSubmitting = true;
				emojis.forEach((e) => form.data.shortcodes.push(e.shortcode));

				const res = await backend.uploadEmojis(form.data);
				res.match(
					(emojis) => {
						userStore.emojis = [...userStore.emojis, ...emojis];
					},
					(err) => logErr(err)
				);

				for (const emoji of emojis) {
					URL.revokeObjectURL(emoji.url);
				}

				emojis = [];
				isSubmitting = false;
			}
		}
	});

	function onFile(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files) {
			for (const image of target.files) {
				const dataUrl = URL.createObjectURL(image);

				emojis.push({
					id: createId(),
					url: dataUrl,
					shortcode: transformShortcode(image.name.split('.')[0])
				});

				$form.emojis = [...$form.emojis, image];
			}
		}
	}

	async function onExistingEmojiDelete(id: string) {
		let deleteTimeout = setTimeout(() => {
			isDeleting = true;
		}, 200);

		const res = await backend.deleteEmoji(id);
		res.match(
			() => {
				userStore.emojis = userStore.emojis.filter((emoji) => emoji.id !== id);
			},
			(err) => logErr(err)
		);

		clearTimeout(deleteTimeout);
		isDeleting = false;
	}

	function onEmojiDelete(id: string, idx?: number) {
		emojis = emojis.filter((emoji) => emoji.id !== id);
		$form.emojis = $form.emojis.filter((_, index) => index !== idx);
	}
</script>

<p class="text-main-400 mt-3">Add emojis to express yourself!</p>

<h2 class="text-main-400 mt-4 text-sm uppercase">Requirements</h2>
<ul class="text-main-400 mt-1 flex flex-col">
	<li>- File type: JPEG, PNG, GIF, WEBP, AVIF</li>
	<li>- Recommended emoji dimensions: 128x128</li>
</ul>

<div class="mt-4 flex items-end gap-x-3">
	<label
		for="avatar-profile"
		class="group bg-accent-darker border-[0.5px] hover:bg-accent border-accent relative flex w-fit items-center justify-center overflow-hidden px-2.5 py-1 transition duration-100 rounded-sm"
	>
		<input
			type="file"
			id="avatar-profile"
			name="avatar-profile"
			aria-label="Profile avatar and banner"
			class="absolute h-full w-full text-transparent hover:cursor-pointer"
			accept="image/png, image/jpeg, image/gif, image/webp, image/avif"
			onchange={onFile}
			multiple
		/>
		<p>Upload an emoji</p>
	</label>
	{#if isDeleting}
		<p class="text-accent-50" transition:fly={{ duration: 50, y: 5 }}>
			Deleting an existing emoji...
		</p>
	{/if}
</div>

<hr class="mt-5 w-full border-none" style="height: 1px; background-color: var(--color-main-800);" />

{#if emojis.length > 0 || userStore.emojis?.length > 0}
	<form use:enhance class="flex flex-col gap-y-3">
		<ul>
			{#each userStore.emojis as emoji (emoji.id)}
				<EmojiLine
					{...emoji}
					bind:shortcode={emoji.shortcode}
					deleteFunction={onExistingEmojiDelete}
				/>
			{/each}
			{#each emojis as emoji, idx (emoji.id)}
				<EmojiLine
					{...emoji}
					{idx}
					bind:shortcode={emoji.shortcode}
					deleteFunction={onEmojiDelete}
				/>
			{/each}
		</ul>

		{#if emojis.length > 0}
			<SubmitButton
				bind:isSubmitting
				bind:isComplete={isButtonComplete}
				text="Save my emojis"
				type="submit"
				class="py-1!"
			/>
		{/if}
	</form>
{/if}

{#if $errors.emojis || $errors.shortcodes}
	<p class="mt-5 flex items-center gap-x-2 text-red-400">
		<DangerIcon height={20} width={20} />
		{$errors.emojis?.[0] || $errors.shortcodes?.[0]}
	</p>
{/if}
