<script lang="ts">
	import { generateHTML } from '@tiptap/core';
	import { MESSAGE_EXTENSIONS } from 'ui/RichInput/richInputConfig';
	import type { User } from '$lib/types/types';

	interface Props {
		user: User;
	}

	let { user }: Props = $props();
</script>

{#if user.about_me}
	<p class="text-white/65 relative z-[4] mix-blend-plus-lighter mt-3 text-sm leading-4.5">
		{@html generateHTML(user.about_me, MESSAGE_EXTENSIONS)}
	</p>
{/if}
