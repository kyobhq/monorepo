<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import ContextMenuUser from 'ui/ContextMenu/ContextMenuUser.svelte';
	import AnimatedAvatar from 'ui/AnimatedAvatar/AnimatedAvatar.svelte';

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
	let hoverAvatar = $state(false);

	function getHoverColor(color: string) {
		return color + '40';
	}
</script>

<button
	bind:this={memberEl}
	style="--hover-color: {color ? getHoverColor(color) : 'var(--ui-main-900)'}"
	class={[
		'relative flex items-center gap-x-2.5 p-1 pr-2.5 transition-colors duration-100 select-none rounded-sm z-[1] active-scale-down',
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
	onmouseenter={() => (hoverAvatar = true)}
	onmouseleave={() => (hoverAvatar = false)}
>
	<div
		class={[
			'relative w-8 h-8 rounded-[2px] overflow-hidden',
			border && 'after:absolute after:inset-0 after:inner-main-700'
		]}
	>
		<AnimatedAvatar src={avatar} alt={name} class="w-full h-full" hover={hoverAvatar} />
	</div>
	<span style="color: {color ? color : 'var(--ui-main-300)'};">{name}</span>

	<ContextMenuUser memberID={id} />
</button>
