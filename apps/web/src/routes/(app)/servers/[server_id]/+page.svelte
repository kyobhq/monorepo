<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { channelStore } from 'stores/channelStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import StraightFaceEmoji from 'ui/icons/StraightFaceEmoji.svelte';

	const currentServer = $derived(serverStore.getServer(page.params.server_id!));

	$effect(() => {
		if (!page.params.server_id) return;
		if (page.params.channel_id) return;

		const firstChannelID = channelStore.getFirstChannel(page.params.server_id);
		if (firstChannelID) goto(`/servers/${page.params.server_id}/channels/${firstChannelID}`);
	});
</script>

{#if currentServer}
	<div
		class="h-full w-full flex items-center justify-center text-2xl flex-col gap-y-4 text-main-700"
	>
		<StraightFaceEmoji height={128} width={128} />
		<p>You're not in any channels yet</p>
	</div>
{/if}
