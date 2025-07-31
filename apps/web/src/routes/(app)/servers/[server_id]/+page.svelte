<script>
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { channelStore } from 'stores/channelStore.svelte';

	$effect(() => {
		if (!page.params.server_id) return;
		if (page.params.channel_id) return;

		const firstChannelID = channelStore.getFirstChannel(page.params.server_id);
		if (firstChannelID) goto(`/servers/${page.params.server_id}/channels/${firstChannelID}`);
	});
</script>

<h1>Server {page.params.server_id}</h1>
