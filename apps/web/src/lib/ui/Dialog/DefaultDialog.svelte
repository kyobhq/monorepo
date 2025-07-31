<script lang="ts">
	import { Dialog } from 'bits-ui';
	import CloseIcon from 'ui/icons/CloseIcon.svelte';
	import { scaleBlur } from 'utils/transition';

	let { children, title, subtitle, state = $bindable() } = $props();
</script>

<Dialog.Root onOpenChange={(s) => (state = s)} open={state}>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20 transition-opacity" />
		<Dialog.Content forceMount={true}>
			{#snippet child({ props, open })}
				{#if open}
					<div
						class="bg-main-950 border-[0.5px] border-main-700 fixed top-1/2 left-1/2 z-50 w-[550px] -translate-1/2"
						{...props}
						transition:scaleBlur={{}}
					>
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
				{/if}
			{/snippet}
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
