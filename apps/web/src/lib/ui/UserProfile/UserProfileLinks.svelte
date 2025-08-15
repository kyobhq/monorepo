<script lang="ts">
	import ExternalLinkIcon from 'ui/icons/ExternalLinkIcon.svelte';
	import type { User } from '$lib/types/types';

	interface Props {
		user: User;
	}

	let { user }: Props = $props();
</script>

{#if user.links.length! > 0}
	<p class="text-white/45 mix-blend-plus-lighter text-sm font-medium relative z-[4]">Links</p>
	<div class="flex flex-col gap-y-1 mt-1.5 relative z-[4]">
		{#each user.links as link, idx (idx)}
			<a
				href={link.url}
				class="bg-white/20 mix-blend-plus-lighter border border-white/45 py-3 px-3 flex items-center justify-between hover:bg-white/30 hover:border-white/65 transition-colors duration-100 rounded-md"
			>
				<div class="flex flex-col gap-y-1.5">
					<p class="text-white/65 mix-blend-plus-lighter leading-none font-medium">
						{link.label}
					</p>
					<p class="text-white/45 mix-blend-plus-lighter text-sm leading-none">
						{link.url}
					</p>
				</div>
				<ExternalLinkIcon height={20} width={20} />
			</a>
		{/each}
	</div>
{/if}
