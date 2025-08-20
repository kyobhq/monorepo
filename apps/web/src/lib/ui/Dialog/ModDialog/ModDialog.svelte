<script lang="ts">
	import { Dialog } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import CloseIcon from 'ui/icons/CloseIcon.svelte';
	import { scaleBlur } from 'utils/transition';

	interface Props {
		title: string;
		subtitle?: string;
		children: Snippet;
		open: boolean;
		openState: any;
	}

	let { title, subtitle, children, open, openState = $bindable() }: Props = $props();
</script>

<Dialog.Root onOpenChange={(s) => (openState = s)} {open}>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/45 z-[998] transition-opacity" />
		<Dialog.Content forceMount={true}>
			{#snippet child({ props, open })}
				{#if open}
					<div
						class="bg-main-950 border-[2px] border-main-850 fixed top-1/2 left-1/2 w-[550px] -translate-1/2 rounded-xl shadow-box z-[999]"
						{...props}
						transition:scaleBlur={{}}
					>
						<div class="relative w-full py-5 px-4">
							<p class="font-semibold text-xl">{title}</p>
							{#if subtitle}
								<p class="text-main-400 text-sm whitespace-pre-line max-w-[28rem]">{subtitle}</p>
							{/if}

							<Dialog.Close
								type="button"
								class="text-main-400 h-8 w-8 hover:bg-main-900 hover:text-main-50 absolute top-3 right-3 transition-colors hover:cursor-pointer border-[0.5px] border-main-700 aspect-square flex justify-center items-center duration-75 rounded-md"
							>
								<CloseIcon height={18} width={18} />
							</Dialog.Close>

							{@render children()}
						</div>
					</div>
				{/if}
			{/snippet}
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
