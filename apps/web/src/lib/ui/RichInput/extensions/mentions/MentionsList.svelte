<script lang="ts">
	import type { SuggestionProps } from '@tiptap/suggestion';
	import { onDestroy } from 'svelte';

	interface Props {
		props: SuggestionProps<any, any>;
		class: string;
	}

	let { props, class: classes }: Props = $props();

	let selectedIndex = $state(0);
	let scrollableMenu = $state<HTMLDivElement>();

	export function handleKeyDown({ event }: { event: KeyboardEvent }) {
		if (props.query === 'everyone') {
			props.command({ 'user-id': 'everyone', label: 'everyone' });
			return true;
		}

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
		if (scrollableMenu) {
			const selectedItem = scrollableMenu.children[selectedIndex] as HTMLElement;
			if (selectedItem) {
				const scrollTop = scrollableMenu.scrollTop;
				const scrollBottom = scrollTop + scrollableMenu.clientHeight;
				const elementTop = selectedItem.offsetTop;
				const elementBottom = elementTop + selectedItem.offsetHeight;

				if (elementTop < scrollTop) {
					scrollableMenu.scrollTop = elementTop;
				} else if (elementBottom > scrollBottom) {
					scrollableMenu.scrollTop = elementBottom - scrollableMenu.clientHeight;
				}
			}
		}
	}

	function handleEnter() {
		selectItem(selectedIndex);
	}

	function selectItem(index: number) {
		const item = props.items[index];

		if (item) {
			props.command({ 'user-id': item.id, label: item.display_name });
		}
	}

	onDestroy(() => {
		selectedIndex = 0;
	});
</script>

{#if props.items.length > 0}
	<div
		bind:this={scrollableMenu}
		class={[
			'bg-main-900 border-[0.5px] border-main-700 z-[10] flex max-h-[20rem] flex-col gap-y-1 overflow-y-auto px-1 py-1',
			classes
		]}
	>
		{#each props.items as item, idx (idx)}
			<button
				class={[
					'flex w-full items-center gap-x-1.5 px-2 py-1 text-left',
					idx === selectedIndex ? 'bg-main-800 text-accent-50' : 'hover:bg-accent-100/20'
				]}
				onclick={() => (selectedIndex = idx)}
			>
				<img src={item.avatar} alt="avatar" class="h-5 w-5 object-cover" />
				{item.display_name}
			</button>
		{/each}
	</div>
{/if}
