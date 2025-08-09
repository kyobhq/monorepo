<script lang="ts">
	import { Switch } from 'bits-ui';

	interface Props {
		active: boolean;
		action: () => void;
		reverse: () => void;
		label?: string;
	}

	let { active, action, reverse, label }: Props = $props();
</script>

<label
	class={[
		'flex items-center gap-x-3 transition-colors duration-75',
		active ? 'text-main-50' : 'text-main-500'
	]}
>
	<Switch.Root
		checked={active}
		onCheckedChange={(checked) => (checked ? action() : reverse())}
		class="data-[state=checked]:bg-accent data-[state=unchecked]:bg-main-800 peer border-main-700 data-[state=checked]:border-accent-lighter inline-flex h-[24px] w-[55px] shrink-0 cursor-pointer items-center border !p-0 px-[3px] transition-colors focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 rounded-sm"
	>
		<Switch.Thumb
			class="bg-main-50 pointer-events-none block size-[20px] shrink-0 transition-transform data-[state=checked]:translate-x-8 data-[state=unchecked]:translate-x-[0.1rem] rounded-[2px]"
		/>
	</Switch.Root>
	<span class="select-none">{label}</span>
</label>
