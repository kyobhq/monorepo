<script lang="ts">
	interface Props {
		avatar: string;
		name: string;
		border?: boolean;
		hoverable?: boolean;
		color?: string;
	}

	const { avatar, name, border = false, hoverable = false, color }: Props = $props();

	function getHoverColor(color: string) {
		return color + '26';
	}
</script>

<div
	style="--hover-color: {getHoverColor(color)}"
	class={[
		'flex items-center gap-x-2.5 p-1 pr-2.5 transition-colors duration-75 select-none',
		hoverable && `hover:cursor-pointer hover:bg-[var(--hover-color)]`
	]}
>
	<div class={['relative w-8 h-8', border && 'after:absolute after:inset-0 after:inner-main-700']}>
		<img src={avatar} alt={name} class="w-full h-full object-cover" />
	</div>
	<span class="text-sm" style="color: {color ? color : 'var(--ui-main-300)'};">{name}</span>
</div>
