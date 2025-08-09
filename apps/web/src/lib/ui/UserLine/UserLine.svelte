<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import { userStore } from 'stores/userStore.svelte';

	interface Props {
		id: string;
		avatar: string;
		name: string;
		border?: boolean;
		hoverable?: boolean;
		color?: string;
		status: string;
	}

	const {
		id,
		avatar = $bindable(),
		name = $bindable(),
		border = false,
		hoverable = false,
		color,
		status
	}: Props = $props();

	let memberEl = $state<HTMLButtonElement>();

	function getHoverColor(color: string) {
		return color + '26';
	}
</script>

<button
	bind:this={memberEl}
	style="--hover-color: {getHoverColor(color ?? '#ADADB8')}"
	class={[
		'flex items-center gap-x-2.5 p-1 pr-2.5 transition-colors duration-75 select-none rounded-sm',
		hoverable && `hover:cursor-pointer hover:bg-[var(--hover-color)]`,
		status === 'offline' && 'opacity-40'
	]}
	onclick={() => {
		if (hoverable && id === userStore.user!.id) {
			coreStore.openMyProfile(memberEl!, 'left');
		} else {
			coreStore.openProfile(id, memberEl!, 'left');
		}
	}}
>
	<div
		class={[
			'relative w-8 h-8 rounded-[2px] overflow-hidden',
			border && 'after:absolute after:inset-0 after:inner-main-700'
		]}
	>
		<img src={avatar} alt={name} class="w-full h-full object-cover" />
	</div>
	<span style="color: {color ? color : 'var(--ui-main-300)'};">{name}</span>
</button>
