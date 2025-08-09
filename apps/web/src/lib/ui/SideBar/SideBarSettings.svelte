<script lang="ts">
	import { Dialog } from 'bits-ui';
	import CloseIcon from 'ui/icons/CloseIcon.svelte';

	interface Props {
		settings: string[];
		navigationFn: (setting: string) => void;
		activeSection: string;
	}

	let { settings, navigationFn, activeSection }: Props = $props();
</script>

<aside
	class="bg-main-950 h-full border-r-[0.5px] border-r-main-700 w-[15rem] overflow-hidden p-2 shrink-0 rounded-l-md"
>
	<Dialog.Close
		type="button"
		class="text-main-400 h-8 w-8 absolute z-50 bg-main-950 top-3 right-3 shrink-0 hocus:bg-main-900 hocus:text-main-50 transition-colors hover:cursor-pointer border-[0.5px] border-main-700 aspect-square flex justify-center items-center duration-75 rounded-sm"
	>
		<CloseIcon height={20} width={20} />
	</Dialog.Close>
	<section class="flex flex-col h-full gap-y-1">
		{#each settings as setting, idx (idx)}
			<button
				class={[
					'text-left px-2.5 py-1.5 transition-colors duration-75 hover:cursor-pointer rounded-sm',
					activeSection === setting
						? 'text-main-50 bg-main-800/60'
						: 'text-main-300 hocus:bg-main-800/80 hocus:text-main-200'
				]}
				onclick={() => navigationFn(setting)}
			>
				{setting}
			</button>
		{/each}
	</section>
</aside>
