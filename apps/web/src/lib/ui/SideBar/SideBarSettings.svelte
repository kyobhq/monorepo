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
	class="bg-main-950 h-full border-r-[2px] border-r-main-850 inner-shadow-box w-[13rem] overflow-hidden p-2 shrink-0 rounded-l-[10px]"
>
	<Dialog.Close
		type="button"
		class="active:scale-[0.95] text-main-400 h-8 w-8 absolute z-50 bg-main-950 top-3 right-3 shrink-0 hover:bg-main-900 hover:text-main-50 transition hover:cursor-pointer border-[0.5px] border-main-700 aspect-square flex justify-center items-center duration-75 rounded-md"
	>
		<CloseIcon height={20} width={20} />
	</Dialog.Close>
	<section class="flex flex-col h-full gap-y-1">
		{#each settings as setting, idx (idx)}
			<button
				class={[
					'text-left px-2.5 py-1.5 transition duration-100 hover:cursor-pointer rounded-lg',
					activeSection === setting
						? 'active-setting'
						: 'text-main-300 hover:bg-main-800/60 hover:text-main-200 active:bg-main-800/80'
				]}
				onclick={() => navigationFn(setting)}
			>
				{setting}
			</button>
		{/each}
	</section>
</aside>
