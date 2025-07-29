<script lang="ts">
	import type { Component } from 'svelte';
	import EyeClose from 'ui/icons/EyeClose.svelte';
	import EyeOpen from 'ui/icons/EyeOpen.svelte';

	interface Props {
		id: string;
		label: string;
		type: 'password' | 'text' | 'email';
		placeholder: string;
		Icon?: Component;
		IWidth?: number;
		formInput: string;
		errorInput?: any;
	}

	let {
		id,
		label,
		type,
		placeholder,
		Icon,
		IWidth,
		formInput = $bindable(),
		errorInput
	}: Props = $props();

	let showPassword = $state(false);
</script>

{#if type !== 'password'}
	<div class="flex flex-col gap-y-1">
		<div class="flex items-center">
			<label for={id} class="text-sm text-main-500">{label}</label>
			{#if errorInput}
				<p class="text-sm leading-none text-red-400">{errorInput}</p>
			{/if}
		</div>

		<div
			class="flex gap-x-2 bg-main-950 w-full border border-main-800 items-center pl-2.5 pr-2 group hocus:bg-main-900 hocus:border-main-700 focus-within:bg-main-900 focus-within:border-main-700 transition-colors duration-75 text-main-300"
		>
			<input
				{type}
				{id}
				bind:value={formInput}
				name={id}
				{placeholder}
				class="w-full py-2 placeholder:text-main-600 text-main-50"
			/>
			<Icon width={IWidth} />
		</div>
	</div>
{:else}
	<div class="flex flex-col gap-y-1">
		<div class="flex items-center">
			<label for={id} class="text-sm text-main-500">{label}</label>
			{#if errorInput}
				<p class="text-sm leading-none text-red-400">{errorInput}</p>
			{/if}
		</div>

		<div
			class="flex gap-x-2 bg-main-950 w-full border border-main-800 items-center pl-2.5 pr-2 group hocus:bg-main-900 hocus:border-main-700 focus-within:bg-main-900 focus-within:border-main-700 transition-colors duration-75 text-main-300"
		>
			<input
				type={showPassword ? 'text' : 'password'}
				{id}
				name={id}
				{placeholder}
				class="w-full py-2 placeholder:text-main-600 text-main-50"
				bind:value={formInput}
			/>
			<button
				type="button"
				onclick={() => (showPassword = !showPassword)}
				class="hover:cursor-pointer text-main-300 hocus:text-main-100 transition-colors"
			>
				{#if showPassword}
					<EyeClose height={20} width={20} />
				{:else}
					<EyeOpen height={20} width={20} />
				{/if}
			</button>
		</div>
	</div>
{/if}
