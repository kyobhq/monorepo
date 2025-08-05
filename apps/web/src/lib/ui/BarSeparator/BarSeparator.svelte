<script lang="ts">
	interface Props {
		title: string;
	}

	const { title }: Props = $props();
</script>

<div class="bg-main-975 border-t-[0.5px] border-b-[0.5px] border-main-700 py-3 px-4 select-none">
	{title}
</div>
