<script lang="ts">
	import { onDestroy } from 'svelte';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';

	interface Props {
		attachments: File[];
	}

	let { attachments = $bindable() }: Props = $props();
	let objectUrls = $state<Record<number, string | null>>([]);

	function removeAttachment(file: File, idx: number) {
		attachments = attachments.filter((f) => f.name !== file.name);
		delete objectUrls[idx];
	}

	function getFileType(file: File): 'image' | 'video' | 'unknown' {
		if (file.type.startsWith('image/')) return 'image';
		if (file.type.startsWith('video/')) return 'video';
		return 'unknown';
	}

	$effect(() => {
		attachments.forEach((file, idx) => {
			const fileType = getFileType(file);
			if (fileType === 'image' || fileType === 'video') {
				objectUrls[idx] = URL.createObjectURL(file);
			} else {
				objectUrls[idx] = null;
			}
		});
	});

	onDestroy(() => {
		Object.values(objectUrls).forEach((url) => url && URL.revokeObjectURL(url));
	});
</script>

{#snippet closeButton(attachment: File, idx: number)}
	<button
		class="hocus:border-red-400 hocus:bg-red-400/40 absolute top-1 right-1 border border-red-400/50 bg-red-400/20 transition-colors duration-100 hover:cursor-pointer z-[2] rounded-sm"
		onclick={() => removeAttachment(attachment, idx)}
	>
		<PlusIcon height={16} width={16} class="text-red-400 rotate-45" />
	</button>
{/snippet}

<div class="flex gap-x-2 p-1 bg-main-900 border-[0.5px] border-main-700 rounded-md items-center">
	{#each attachments as attachment, idx (idx)}
		{@const fileType = getFileType(attachment)}

		{#if fileType === 'image' && objectUrls[idx]}
			<figure class="relative aspect-square h-14 w-14 rounded-sm overflow-hidden highlight-border">
				{@render closeButton(attachment, idx)}
				<img src={objectUrls[idx]} alt={attachment.name} class="h-full w-full object-cover" />
			</figure>
		{:else if fileType === 'video'}
			<figure class="relative aspect-square h-14 w-14 rounded-sm overflow-hidden highlight-border">
				{@render closeButton(attachment, idx)}
				<video class="h-full w-full object-cover">
					<source src={objectUrls[idx]} />
					<track kind="captions" />
				</video>
			</figure>
		{:else}
			<figure
				class="relative flex aspect-square h-14 w-14 flex-col items-center justify-center gap-y-1 overflow-hidden rounded-sm bg-main-800 border-[0.5px] border-main-600"
			>
				{@render closeButton(attachment, idx)}
				<p class="w-[calc(100%-1.5rem)] truncate">{attachment.name}</p>
			</figure>
		{/if}
	{/each}
</div>
