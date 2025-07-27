<script lang="ts">
	interface Props {
		image: string;
		active?: boolean;
		onclick: () => void;
	}

	let { image, active, onclick }: Props = $props();
</script>

<button
	class={[
		'h-13 w-13 aspect-square object-cover relative after:content-normal after:absolute after:inset-0 after:transition',
		active
			? 'after:inner-active'
			: 'after:inner-main-700 hocus:after:inner-active hover:cursor-pointer'
	]}
	{onclick}
>
	<img src={image} alt="server-icon" class="h-full w-full" />
</button>
