<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { coreStore } from 'stores/coreStore.svelte';
	import CloseIcon from 'ui/icons/CloseIcon.svelte';
	import { scaleBlur } from 'utils/transition';
</script>

<Dialog.Root
	onOpenChange={(s) => (coreStore.destructiveDialog.open = s)}
	open={coreStore.destructiveDialog.open}
>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20 transition-opacity" />
		<Dialog.Content forceMount={true}>
			{#snippet child({ props, open })}
				{#if open}
					<div
						class="bg-main-950 border-[0.5px] border-main-700 fixed top-1/2 left-1/2 z-50 w-[550px] -translate-1/2 rounded-md"
						{...props}
						transition:scaleBlur={{}}
					>
						<div class="relative w-full py-5 px-4">
							<p class="font-semibold text-xl">{coreStore.destructiveDialog.title}</p>
							<p class="text-main-400 text-sm">{coreStore.destructiveDialog.subtitle}</p>

							<Dialog.Close
								type="button"
								class="text-main-400 h-8 w-8 hocus:bg-main-900 hocus:text-main-50 absolute top-3 right-3 transition-colors hover:cursor-pointer border-[0.5px] border-main-700 aspect-square flex justify-center items-center duration-75 rounded-sm"
							>
								<CloseIcon height={18} width={18} />
							</Dialog.Close>
						</div>

						<div class="flex justify-end p-3">
							<button
								type="button"
								onclick={coreStore.destructiveDialog.onclick}
								class="bg-red-400/20 border-[0.5px] border-red-400 px-2 py-1 hover:cursor-pointer hocus:bg-red-400 transition-colors duration-75 text-red-400 hocus:text-red-50 rounded-sm"
							>
								{coreStore.destructiveDialog.buttonText}
							</button>
						</div>
					</div>
				{/if}
			{/snippet}
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
