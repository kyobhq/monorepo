<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import BarSeparator from 'ui/BarSeparator/BarSeparator.svelte';
	import Channel from 'ui/Channel/Channel.svelte';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
	import ContextMenuSideBar from 'ui/ContextMenu/ContextMenuSideBar.svelte';

	const currentServer = $derived(serverStore.getServer(page.params.server_id || '') || undefined);
</script>

{#if currentServer}
	<BarSeparator title={currentServer.name} />
	<section class="relative flex flex-col gap-y-2 p-2.5 h-full">
		<ContextMenuSideBar />
		{#each Object.values(currentServer.categories).sort((a, b) => a.position - b.position) as category (category.id)}
			<CollapsibleBox header={category.name} categoryId={category.id}>
				{#if category.channels}
					{#each Object.values(category.channels).sort((a, b) => a.position - b.position) as channel (channel.id)}
						<Channel
							id={channel.id}
							categoryId={category.id}
							type={channel.type}
							name={channel.name}
							onclick={() => goto(`/servers/${channel.server_id}/channels/${channel.id}`)}
							active={page.url.pathname.includes(channel.id)}
						/>
					{/each}
				{/if}
			</CollapsibleBox>
		{/each}
	</section>
{/if}
