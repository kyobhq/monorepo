<script lang="ts">
	import type { SuggestionProps } from '@tiptap/suggestion';
	import { editorStore } from 'stores/editorStore.svelte';

	interface Props {
		props: SuggestionProps<any, any>;
		class: string;
	}

	let { props, class: classes }: Props = $props();

	let selectedIndex = $state(0);

	export function handleKeyDown({ event }: { event: KeyboardEvent }) {
		if (event.key === 'ArrowUp') {
			handleArrowUp();
			return true;
		}

		if (event.key === 'ArrowDown') {
			handleArrowDown();
			return true;
		}

		if (event.key === 'Enter') {
			handleEnter();
			return true;
		}
	}

	function handleArrowUp() {
		selectedIndex = (selectedIndex + props.items.length - 1) % props.items.length;
		scrollToSelectedItem();
	}

	function handleArrowDown() {
		selectedIndex = (selectedIndex + 1) % props.items.length;
		scrollToSelectedItem();
	}

	function scrollToSelectedItem() {
		if (editorStore.menuScrollContainer) {
			const selectedItem = editorStore.menuScrollContainer.children[selectedIndex] as HTMLElement;
			if (selectedItem) {
				const scrollTop = editorStore.menuScrollContainer.scrollTop;
				const scrollBottom = scrollTop + editorStore.menuScrollContainer.clientHeight;
				const elementTop = selectedItem.offsetTop;
				const elementBottom = elementTop + selectedItem.offsetHeight;

				if (elementTop < scrollTop) {
					editorStore.menuScrollContainer.scrollTop = elementTop;
				} else if (elementBottom > scrollBottom) {
					editorStore.menuScrollContainer.scrollTop =
						elementBottom - editorStore.menuScrollContainer.clientHeight;
				}
			}
		}
	}

	function handleEnter() {
		selectItem(selectedIndex);
	}

	function selectItem(index: number) {
		const emoji = props.items[index].emoji;

		if (emoji.url) {
			props.command({ url: emoji.url, label: emoji.shortcode });
		} else {
			props.command({ emoji: emoji.unicode, label: emoji.label });
		}
	}

	$effect(() => {
		if (props.items.length) {
			selectedIndex = 0;
		}
	});
</script>

{#if props.items.length > 0}
	<div
		bind:clientHeight={editorStore.menuHeight}
		class={[
			'bg-main-900 border-[0.5px] border-main-700 z-[10] flex max-h-[20rem] flex-col gap-y-1 overflow-y-auto px-1 py-1',
			classes
		]}
	>
		{#each props.items as item, idx (idx)}
			{@const emoji = item.emoji}
			<button
				class={[
					'flex w-full items-center gap-x-1.5 px-2 py-1 text-left',
					idx === selectedIndex ? 'bg-main-800 text-accent-50' : 'hover:bg-accent-100/20'
				]}
				onclick={() => (selectedIndex = idx)}
			>
				{#if emoji.url}
					<img src={emoji.url} alt={emoji.label} class="h-[24px] w-[24px] object-contain" />
					:{emoji.shortcode}:
				{:else}
					{emoji.unicode} :{emoji.label.replaceAll(' ', '_')}:
				{/if}
			</button>
		{/each}
	</div>
{/if}
