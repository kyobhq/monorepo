<script lang="ts">
	import { onMount } from 'svelte';
	import type { Attachment } from '$lib/types/types';
	import { getFileType } from 'utils/files';
	import AttachmentImage from './attachments/AttachmentImage.svelte';
	import AttachmentVideo from './attachments/AttachmentVideo.svelte';
	import AttachmentFile from './attachments/AttachmentFile.svelte';

	let { attachments, hover }: { attachments: Attachment[]; hover: boolean } = $props();

	let images = $state<Attachment[]>([]);
	let videos = $state<Attachment[]>([]);
	let files = $state<Attachment[]>([]);

	onMount(() => {
		for (const attachment of attachments) {
			const fileType = getFileType(attachment.url);

			switch (fileType) {
				case 'image':
					images.push(attachment);
					break;
				case 'video':
					videos.push(attachment);
					break;
				case 'unknown':
					files.push(attachment);
					break;
			}
		}
	});

	const medias = $derived([...images, ...videos]);
</script>

{#if medias.length > 0}
	<div class={medias.length > 1 ? 'grid w-fit max-w-[50%] grid-cols-3 gap-2' : 'max-w-[15%]'}>
		{#each images as image, idx (idx)}
			<AttachmentImage {image} {medias} {hover} />
		{/each}
		{#each videos as video, idx (idx)}
			<AttachmentVideo {video} />
		{/each}
	</div>
{/if}

{#if files.length > 0}
	<div class="pointer-events-auto flex flex-col">
		{#each files as file, idx (idx)}
			<AttachmentFile {file} />
		{/each}
	</div>
{/if}

<style>
	:global(.attachment::before) {
		content: '';
		position: absolute;
		inset: 0;
		box-shadow: inset 0 0 0 1px #fafafa33;
		transition: box-shadow ease-out 75ms;
	}

	:global(.attachment:hover::before) {
		box-shadow: inset 0 0 0 1px #fafafa;
		cursor: pointer;
	}
</style>
