<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { scaleBlur } from 'utils/transition';

	let { children, state = $bindable() } = $props();
</script>

<Dialog.Root onOpenChange={(s) => (state = s)} open={state}>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20 transition-opacity" />
		<Dialog.Content forceMount={true}>
			{#snippet child({ props, open })}
				{#if open}
					<div
						class="bg-main-950 border-[0.5px] border-main-700 fixed top-1/2 left-1/2 z-50 w-[70rem] h-[50rem] -translate-1/2 flex"
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
