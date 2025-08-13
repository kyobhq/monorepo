<script lang="ts">
	import FormSwitch from 'ui/Form/FormSwitch.svelte';
	import { ABILITIES_CONTENT } from '$lib/constants/permissions';

	let { form = $bindable() } = $props();
</script>

<div>
	{#each ABILITIES_CONTENT as ability, idx (idx)}
		<FormSwitch
			active={form.abilities.includes(ability.code)}
			action={() => (form.abilities = [...form.abilities, ability.code])}
			reverse={() =>
				(form.abilities = form.abilities.filter((code: string) => code !== ability.code))}
			class={idx === ABILITIES_CONTENT.length - 1 ? 'border-b-[0px]' : 'mt-2'}
			label={ability.label}
			description={ability.description}
		/>
	{/each}
</div>
