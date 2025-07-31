<script lang="ts">
	import { fly } from 'svelte/transition';
	import Check from 'ui/icons/Check.svelte';
	import LoadingIcon from 'ui/icons/LoadingIcon.svelte';

	interface Props {
		type?: 'button' | 'submit';
		text: string;
		isSubmitting: boolean;
		isComplete: boolean;
		isEmpty?: boolean;
	}

	let {
		type = 'submit',
		text,
		isSubmitting = $bindable(),
		isComplete = $bindable(),
		isEmpty = $bindable()
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
			}, 1000);
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
		'px-3 py-1.5 bg-accent-darker border-[0.5px] border-accent hocus:bg-accent transition-colors hover:cursor-pointer relative'
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
