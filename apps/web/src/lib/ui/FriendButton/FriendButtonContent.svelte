<script lang="ts">
	import { generateHTML } from '@tiptap/core';
	import { MESSAGE_EXTENSIONS } from 'ui/RichInput/richInputConfig';

	interface Props {
		avatar?: string;
		display_name?: string;
		about_me?: any;
		accepted: boolean;
		sender: boolean;
	}

	let { avatar, display_name, about_me, accepted, sender }: Props = $props();
</script>

<figure class="h-10 w-10 rounded-lg overflow-hidden">
	<img src={avatar} class="h-full w-full object-cover" alt="" />
</figure>
<div class="flex flex-col">
	<p class="font-medium">
		{display_name}
		{#if !accepted && sender}
			<span class="text-sm text-main-500">- Pending...</span>
		{/if}
	</p>
	{#if about_me}
		<div class="text-sm text-main-500 w-[8rem] bio-truncate">
			{@html generateHTML(about_me, MESSAGE_EXTENSIONS)}
		</div>
	{/if}
</div>

<style>
	.bio-truncate {
		display: -webkit-box;
		line-clamp: 1;
		-webkit-line-clamp: 1;
		-webkit-box-orient: vertical;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.bio-truncate :global(*) {
		display: inline !important;
		line-height: inherit;
	}
</style>
