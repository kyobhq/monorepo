<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { scaleBlur } from 'utils/transition';

	let { children, state = $bindable(), initialized = $bindable() } = $props();

	$effect(() => {
		if (state) {
			setTimeout(() => {
				initialized = true;
			}, 500);
		} else {
			initialized = false;
		}
	});
</script>

<Dialog.Root onOpenChange={(s) => (state = s)} open={state}>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/45 transition-opacity z-[998]" />
		<Dialog.Content forceMount={true}>
			{#snippet child({ props, open })}
				{#if open}
					<div
						class="bg-main-950 border-[2px] border-main-850 shadow-box fixed top-1/2 left-1/2 z-[999] w-[60rem] h-[40rem] -translate-1/2 flex rounded-xl"
						{...props}
						transition:scaleBlur={{}}
					>
						{@render children()}
					</div>
				{/if}
			{/snippet}
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
