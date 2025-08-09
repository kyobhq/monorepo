<script lang="ts">
	import PlusIcon from 'ui/icons/PlusIcon.svelte';
	import { createId } from '@paralleldrive/cuid2';

	interface Props {
		links: { id: string; label: string; url: string }[];
	}

	let { links = $bindable() }: Props = $props();

	const MAX_LINKS = 2;

	function addLink() {
		links = [
			...links,
			{
				id: createId(),
				label: '',
				url: ''
			}
		];
	}

	function removeLink(id: string) {
		links = links.filter((link) => link.id !== id);
	}

	function updateLinkLabel(id: string, newLabel: string) {
		links = links.map((link) => (link.id === id ? { ...link, label: newLabel } : link));
	}

	function updateLinkUrl(id: string, newUrl: string) {
		links = links.map((link) => (link.id === id ? { ...link, url: newUrl } : link));
	}
</script>

<div>
	<p class="text-main-500 text-sm select-none">Links</p>
	<div
		class="flex flex-col bg-main-900 border-[0.5px] border-main-800 placeholder:text-main-400 p-1.5 mt-1.5 transition-colors duration-100 focus:ring-0 rounded-md"
	>
		<ul class="flex flex-col">
			{#each links as link, index (index)}
				<li
					class={['w-full flex gap-x-1 mb-1', links.length < MAX_LINKS ? 'last:mb-2' : 'last:mb-0']}
				>
					<input
						class="bg-main-950 border-[0.5px] border-main-700 w-full px-2 py-1 placeholder:text-main-600 rounded-r-[2px] rounded-l-sm"
						placeholder="My Portfolio"
						value={link.label}
						oninput={(e) => updateLinkLabel(link.id, e.currentTarget.value)}
					/>
					<input
						class="bg-main-950 border-[0.5px] border-main-700 w-full px-2 py-1 placeholder:text-main-600 rounded-[2px]"
						placeholder="https://google.com"
						value={link.url}
						type="url"
						oninput={(e) => updateLinkUrl(link.id, e.currentTarget.value)}
					/>
					<button
						type="button"
						class="bg-red-500/20 border-[0.5px] border-red-400/50 text-red-400 hover:bg-red-500/30 transition-colors duration-100 shrink-0 aspect-square flex items-center justify-center hover:cursor-pointer rounded-r-sm rounded-l-[2px]"
						onclick={() => removeLink(link.id)}
					>
						<PlusIcon height={18} width={18} class="rotate-45" />
					</button>
				</li>
			{/each}
		</ul>

		{#if links.length < MAX_LINKS}
			<button
				type="button"
				class="flex justify-center bg-main-800/40 border-[0.5px] border-main-700/50 placeholder:text-main-400 px-3 py-2 transition-colors duration-100 focus:ring-0 text-main-500 hover:cursor-pointer hocus:bg-main-800 hocus:border-main-700 rounded-sm"
				onclick={addLink}
			>
				<PlusIcon height={20} width={20} />
			</button>
		{/if}
	</div>
</div>
