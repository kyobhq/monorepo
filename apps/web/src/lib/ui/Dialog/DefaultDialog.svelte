<script lang="ts">
	import { Dialog } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import CloseIcon from 'ui/icons/CloseIcon.svelte';
	import { scaleBlur } from 'utils/transition';

	interface Props {
		children: Snippet;
		title: string;
		subtitle?: string;
		state: boolean;
		tabs?: string[];
		currentTab?: string;
	}

	let {
		children,
		title,
		subtitle,
		state = $bindable(),
		tabs,
		currentTab = $bindable()
	}: Props = $props();
</script>

<Dialog.Root
	onOpenChange={(s) => {
		state = s;

		if (tabs && !s) currentTab = tabs[0];
	}}
	open={state}
>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/45 z-[998] transition-opacity" />
		<Dialog.Content forceMount={true}>
			{#snippet child({ props, open })}
				{#if open}
					<div
						class={[
							'fixed top-1/2 left-1/2 z-[999] w-[550px] flex flex-col gap-y-1',
							tabs ? ' -translate-x-1/2 top-[20%]' : '-translate-1/2'
						]}
						{...props}
						transition:scaleBlur={{}}
					>
						{#if tabs}
							<div class="flex gap-x-1">
								{#each tabs as tab, idx (idx)}
									<button
										onclick={() => (currentTab = tab)}
										class={[
											'border-[0.5px] px-2.5 py-1.5 w-fit transition-colors duration-100 hover:cursor-pointer',
											currentTab === tab
												? 'bg-main-800/70 border-main-500'
												: 'bg-main-950 border-main-700 hocus:bg-main-900 hocus:border-main-600'
										]}
									>
										{tab}
									</button>
								{/each}
							</div>
						{/if}
						<div class="bg-main-950 border-[0.5px] border-main-700">
							<div class="border-b-main-800 relative mb-8 w-full border-b py-5 px-4">
								<p class="font-semibold text-xl">{title}</p>
								<p class="text-main-400 text-sm">{subtitle}</p>

								<Dialog.Close
									type="button"
									class="text-main-400 h-8 w-8 hocus:bg-main-900 hocus:text-main-50 absolute top-5 right-5 transition-colors hover:cursor-pointer border-[0.5px] border-main-700 aspect-square flex justify-center items-center duration-75"
								>
									<CloseIcon height={18} width={18} />
								</Dialog.Close>
							</div>
							{@render children?.()}
						</div>
					</div>
				{/if}
			{/snippet}
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
