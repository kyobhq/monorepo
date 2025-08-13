<script lang="ts">
	import ColorPicker from 'ui/ColorPicker/ColorPicker.svelte';
	import RichFormInput from './RichFormInput.svelte';

	interface Props {
		id: string;
		inputValue: any;
		error?: any;
		placeholder: string;
		title: string;
		type: 'text' | 'password' | 'textarea' | 'rich' | 'color-picker' | 'email';
		class?: string;
		inputClass?: string;
	}

	let {
		id,
		inputValue = $bindable(),
		error = $bindable(),
		placeholder,
		title,
		type,
		class: classes,
		inputClass
	}: Props = $props();
</script>

<div class={['flex flex-col', classes]}>
	<div class="flex items-center gap-x-1">
		<label for={id} class={['text-sm select-none mb-1.5', error ? 'text-red-400' : 'text-main-500']}
			>{title}</label
		>
		{#if error}
			<p class="text-sm text-red-400">- {error}</p>
		{/if}
	</div>
	{#if type === 'textarea'}
		<textarea
			{id}
			bind:value={inputValue}
			{placeholder}
			class={[
				'bg-main-900 placeholder:text-main-400 hover:bg-main-800/50 h-30 resize-none transition-colors duration-100 focus:ring-0',
				error ? 'border-red-400' : 'border-main-800 hover:border-main-700'
			]}
		></textarea>
	{:else if type === 'rich'}
		<RichFormInput bind:content={inputValue} {placeholder} class={inputClass} />
	{:else if type === 'color-picker'}
		<ColorPicker bind:color={inputValue} />
	{:else}
		<input
			{id}
			{type}
			bind:value={inputValue}
			{placeholder}
			class={[
				'bg-main-900 border-[0.5px] border-main-800 placeholder:text-main-400 hover:bg-main-800 hover:border-main-600 px-3 py-2 transition-colors duration-100 focus:ring-0 rounded-sm',
				error ? 'border-red-400' : 'border-main-800 hover:border-main-700',
				inputClass
			]}
		/>
	{/if}
</div>
