<script>
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import Channel from 'ui/Channel/Channel.svelte';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
</script>

{#if userStore.pinned_channels.length > 0}
	<div class="px-2.5 pb-2.5">
		<CollapsibleBox header="Pinned channels">
			{#each userStore.pinned_channels as channel (channel.id)}
				{@const channelHref = `/servers/${channel.server_id}/channels/${channel.id}`}
				{@const channelServer = serverStore.getServer(channel.server_id)}

				<Channel
					id={channel.id}
					type={channel.type}
					categoryId={channel.category_id}
					name={channel.name}
					serverName={channelServer.name}
					onclick={() => goto(channelHref)}
					active={page.url.pathname.includes(channelHref)}
				/>
			{/each}
		</CollapsibleBox>
	</div>
{/if}
