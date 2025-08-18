<script lang="ts">
	interface Props {
		avatar?: string;
		display_name?: string;
		about_me?: string;
		sender: boolean;
	}

	let { avatar, display_name, about_me, sender }: Props = $props();
</script>

<figure class="h-10 w-10 rounded-lg overflow-hidden">
	<img src={avatar} class="h-full w-full object-cover" alt="" />
</figure>
<p class="font-medium">
	{display_name}
	{#if sender}
		<span class="text-sm text-main-500">- Pending...</span>
	{/if}
</p>
{#if about_me}
	<p class="text-sm text-main-500">{about_me}</p>
{/if}
