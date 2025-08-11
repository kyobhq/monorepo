<script lang="ts">
	import { backend } from 'stores/backendStore.svelte';
	import { onMount } from 'svelte';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';
	import { transformShortcode } from 'utils/emojis';
	import { logErr } from 'utils/print';

	interface Props {
		id: string;
		idx?: number;
		url: string;
		shortcode: string;
		deleteFunction: (id: string, idx?: number) => void;
	}

	let { id, idx, url, shortcode = $bindable(), deleteFunction }: Props = $props();
	let input = $state<HTMLInputElement>();

	async function onBlurUpdateShortcode(e: any) {
		const target = e.target as HTMLInputElement;

		if (shortcode !== '' && target.value !== shortcode) {
			const res = await backend.updateEmoji(id, target.value);
			res.match(
				() => {
					shortcode = target.value;
				},
				(error) => logErr(error)
			);
		} else {
			shortcode = target.value;
		}
	}

	function onInput(e: any) {
		const target = e.target as HTMLInputElement;
		const transformed = transformShortcode(target.value);
		target.value = transformed;

		const placeholder = target.getAttribute('placeholder') || '';

		const textToMeasure = target.value || placeholder;

		const measurer = document.createElement('span');
		measurer.textContent = textToMeasure;
		measurer.style.cssText = `
        position: absolute;
        visibility: hidden;
        height: auto;
        width: auto;
        white-space: nowrap;
        padding: ${window.getComputedStyle(target).padding};
    `;

		document.body.appendChild(measurer);
		const width = measurer.offsetWidth;
		document.body.removeChild(measurer);

		const minWidth = target.value ? width : Math.max(width, 95);
		target.style.width = `${minWidth}px`;
	}

	onMount(() => {
		onInput({ target: input });
	});
</script>

<li
	class="group hocus:bg-main-800 hocus:border-main-700 focus-within:bg-main-800 focus-within:border-main-700 border-b-main-800 flex items-center justify-between border border-t-transparent border-r-transparent border-l-transparent p-4 transition-colors"
>
	<div class="flex items-center gap-x-2">
		<img src={url} alt={shortcode} class="h-[32px] w-[32px] object-contain" />
		<div
			class="hover:bg-main-900 hover:border-main-700 focus-within:bg-main-900 focus-within:border-main-700 flex h-[32px] items-center border border-transparent px-2 transition-colors rounded-sm"
		>
			<span class="text-main-50/25 font-bold">:</span>
			<input
				type="text"
				value={shortcode}
				bind:this={input}
				oninput={onInput}
				onblur={onBlurUpdateShortcode}
				class="placeholder:text-main-50/50 h-full border-none bg-transparent px-1 py-2 focus-visible:ring-0 focus-visible:outline-none"
				placeholder="emoji_name"
			/>
			<span class="text-main-50/25 font-bold">:</span>
		</div>
	</div>

	<button
		type="button"
		class="inner-red-400/20 hocus:inner-red-400/40 hocus:bg-red-400/25 group-hocus:opacity-100 bg-red-400/15 p-1 text-red-400 opacity-0 transition hover:cursor-pointer rounded-sm"
		aria-label="Delete emoji"
		onclick={() => deleteFunction(id, idx)}
	>
		<PlusIcon height={16} width={16} class="rotate-45" />
	</button>
</li>
