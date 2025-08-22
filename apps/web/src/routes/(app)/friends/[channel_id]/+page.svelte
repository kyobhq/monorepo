<script lang="ts">
	import { page } from '$app/state';
	import { userStore } from 'stores/userStore.svelte';
	import ChatBox from 'ui/ChatBox/ChatBox.svelte';
	import RichInput from 'ui/RichInput/RichInput.svelte';

	const currentFriend = $derived(
		userStore.friends.find((friend) => friend.channel_id === page.params.channel_id)
	);
</script>

{#if currentFriend}
	<ChatBox serverID="global" channelID={currentFriend.channel_id!} />
	<RichInput friend={currentFriend} />
{/if}
