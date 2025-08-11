<script lang="ts">
	import { fly } from 'svelte/transition';
	import Check from 'ui/icons/Check.svelte';
	import LoadingIcon from 'ui/icons/LoadingIcon.svelte';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';

	interface Props {
		type?: 'button' | 'submit';
		text: string;
		isSubmitting: boolean;
		isComplete: boolean;
		isEmpty?: boolean;
		isError?: boolean;
		class?: string;
	}

	let {
		type = 'submit',
		text,
		isSubmitting = $bindable(),
		isComplete = $bindable(),
		isEmpty = $bindable(),
		isError = $bindable(),
		class: classes
	}: Props = $props();

	let wasSubmitting = $state(false);
	let isSubmitted = $state(false);

	$effect(() => {
		isComplete = !isSubmitting && !wasSubmitting && !isSubmitted;
	});

	$effect(() => {
		if (wasSubmitting && !isSubmitting && !isSubmitted) {
			const timer = setTimeout(() => {
				isSubmitted = true;
			}, 600);
			return () => clearTimeout(timer);
		}
		wasSubmitting = isSubmitting;
	});

	$effect(() => {
		if (isSubmitted) {
			const timer = setTimeout(() => {
				isSubmitted = false;
			}, 800);
			return () => clearTimeout(timer);
		}
	});
</script>

<button
	{type}
	class={[
		'px-3 py-1.5 bg-accent-darker border-[0.5px]  transition-colors hover:cursor-pointer relative rounded-sm w-fit',
		isError ? 'bg-red-500/20  border-red-400 hocus:bg-red-400/30' : 'border-accent hocus:bg-accent',
		classes
	]}
	disabled={isSubmitted || isSubmitting || isEmpty}
>
	{#if isSubmitting || wasSubmitting}
		<div
			class="absolute inset-0 flex justify-center items-center"
			transition:fly={{ duration: 200, delay: 150, y: 5 }}
		>
			<LoadingIcon height={20} width={20} />
		</div>
	{:else if isSubmitted && isError}
		<div
			class="absolute inset-0 flex justify-center items-center text-red-400"
			in:fly={{ duration: 200, delay: 200, y: 5 }}
			out:fly={{ duration: 200, y: 5 }}
		>
			<PlusIcon height={20} width={20} class="rotate-45" />
		</div>
	{:else if isSubmitted}
		<div
			class="absolute inset-0 flex justify-center items-center"
			in:fly={{ duration: 200, delay: 200, y: 5 }}
			out:fly={{ duration: 200, y: 5 }}
		>
			<Check height={20} width={20} />
		</div>
	{/if}

	<span
		class={[
			'transition block',
			!isSubmitting && !wasSubmitting && !isSubmitted
				? 'opacity-100 translate-y-0'
				: 'opacity-0 translate-y-1'
		]}
	>
		{text}
	</span>
</button>
