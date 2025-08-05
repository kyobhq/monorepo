<script lang="ts">
	import PlusIcon from 'ui/icons/PlusIcon.svelte';
	import { createId } from '@paralleldrive/cuid2';
	import { Select } from 'bits-ui';

	interface Props {
		facts: { id: string; label: string; value: string }[];
	}

	let { facts = $bindable() }: Props = $props();

	const MAX_FACTS = 3;

	const factBeginnings = [
		{ value: 'I like', label: 'I like' },
		{ value: 'I love', label: 'I love' },
		{ value: 'My birthday is on the', label: 'My birthday is on the' },
		{ value: 'I am from', label: 'I am from' },
		{ value: 'My favorite', label: 'My favorite' },
		{ value: 'I enjoy', label: 'I enjoy' },
		{ value: 'I collect', label: 'I collect' },
		{ value: 'I speak', label: 'I speak' },
		{ value: 'I studied', label: 'I studied' },
		{ value: 'I work as', label: 'I work as' }
	];

	function addFact() {
		facts = [
			...facts,
			{
				id: createId(),
				label: 'I like',
				value: ''
			}
		];
	}

	function removeFact(id: string) {
		facts = facts.filter((fact) => fact.id !== id);
	}

	function updateFactLabel(id: string, newLabel: string) {
		facts = facts.map((fact) => (fact.id === id ? { ...fact, label: newLabel } : fact));
	}

	function updateFactValue(id: string, newValue: string) {
		facts = facts.map((fact) => (fact.id === id ? { ...fact, value: newValue } : fact));
	}
</script>

<div>
	<p class="text-main-500 text-sm select-none">Facts</p>
	<div
		class="flex flex-col bg-main-900 border-[0.5px] border-main-800 placeholder:text-main-400 p-2 mt-1.5 transition-colors duration-100 focus:ring-0"
	>
		<ul class="flex flex-col">
			{#each facts as fact, index (index)}
				<li class="w-full flex gap-x-1 mb-2 last:mb-2">
					<Select.Root
						type="single"
						value={fact.label}
						onValueChange={(newValue) => newValue && updateFactLabel(fact.id, newValue)}
					>
						<Select.Trigger
							class="bg-main-950 border-[0.5px] border-main-700 w-full px-2 py-1 hover:bg-main-800/50 transition-colors duration-100 text-left hover:cursor-pointer"
						>
							{fact.label || 'Select...'}
						</Select.Trigger>
						<Select.Content
							side="top"
							sideOffset={5}
							class="bg-main-950 border-[0.5px] border-main-700 max-h-24 overflow-y-auto z-50 w-[var(--bits-select-anchor-width)] min-w-[var(--bits-select-anchor-width)] p-1"
						>
							{#each factBeginnings as beginning (beginning.value)}
								<Select.Item
									value={beginning.value}
									label={beginning.label}
									class="px-2 py-1 hocus:bg-main-800 cursor-pointer text-main-200 transition-colors duration-100"
								>
									{beginning.label}
								</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
					<input
						class="bg-main-950 border-[0.5px] border-main-700 w-full px-2 py-1 placeholder:text-main-600"
						placeholder="Potatoes"
						value={fact.value}
						oninput={(e) => updateFactValue(fact.id, e.currentTarget.value)}
					/>
					<button
						type="button"
						class="bg-red-500/20 border-[0.5px] border-red-400/50 text-red-400 hover:bg-red-500/30 transition-colors duration-100 shrink-0 aspect-square flex items-center justify-center hover:cursor-pointer"
						onclick={() => removeFact(fact.id)}
					>
						<PlusIcon height={18} width={18} class="rotate-45" />
					</button>
				</li>
			{/each}
		</ul>

		{#if facts.length < MAX_FACTS}
			<button
				type="button"
				class="flex justify-center bg-main-800/40 border-[0.5px] border-main-700/50 placeholder:text-main-400 px-3 py-2 transition-colors duration-100 focus:ring-0 text-main-500 hover:cursor-pointer hocus:bg-main-800 hocus:border-main-700"
				onclick={addFact}
			>
				<PlusIcon height={20} width={20} />
			</button>
		{/if}
	</div>
</div>

