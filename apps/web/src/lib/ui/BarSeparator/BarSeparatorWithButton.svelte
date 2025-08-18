<script lang="ts">
	interface Props {
		title: string;
		buttonText: string;
		onclick: () => void;
	}

	const { title, buttonText, onclick }: Props = $props();
</script>

<div
	class="bg-main-975 border-t-[0.5px] border-b-[0.5px] border-main-800 py-3 px-4 select-none flex items-center justify-between"
>
	{title}
	<button
		{onclick}
		class="text-sm bg-accent-darker px-1.5 py-0.5 hover:bg-accent transition-colors duration-100 rounded-sm"
	>
		{buttonText}
	</button>
</div>
