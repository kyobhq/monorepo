<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import BarSeparator from 'ui/BarSeparator/BarSeparator.svelte';
	import Channel from 'ui/Channel/Channel.svelte';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
	import ServersSlider from 'ui/ServersSlider/ServersSlider.svelte';

	const currentServer = $derived(serverStore.getServer(page.params.server_id || '') || undefined);
</script>

<section class="mt-2.5">
	<ServersSlider />
	{#if userStore.pinned_channels.length > 0}
		<div class="px-2.5">
			<CollapsibleBox header="Pinned channels">
				{#each userStore.pinned_channels as channel (channel.id)}
					{@const channelHref = `/servers/${channel.server_id}/channels/${channel.id}`}
					{@const channelServer = serverStore.getServer(channel.server_id)}

					<Channel
						type={channel.type}
						name={channel.name}
						serverName={channelServer.name}
						onclick={() => goto(channelHref)}
						active={page.url.pathname.includes(channelHref)}
					/>
				{/each}
			</CollapsibleBox>
		</div>
	{/if}
	{#if currentServer}
		<BarSeparator title={currentServer.name} />
		<section class="flex flex-col gap-y-2 p-2.5">
			<CollapsibleBox header="general">
				<Channel type="textual" name="General" onclick={() => {}} />
			</CollapsibleBox>
			<CollapsibleBox header="cool stuff">
				<Channel type="textual" name="General" onclick={() => {}} />
			</CollapsibleBox>
			<CollapsibleBox header="vocals">
				<Channel type="voice" name="General" onclick={() => {}} />
				<Channel type="voice" name="Cowork" onclick={() => {}} />
			</CollapsibleBox>
		</section>
	{/if}
</section>
