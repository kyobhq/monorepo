<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Friend } from '$lib/types/types';
	import { userStore } from 'stores/userStore.svelte';
	import FriendButtonContent from './FriendButtonContent.svelte';
	import FriendshipControls from './FriendshipControls.svelte';
	import { page } from '$app/state';
	import ContextMenuFriend from 'ui/ContextMenu/ContextMenuFriend.svelte';

	interface Props {
		friend: Friend;
	}

	let { friend }: Props = $props();
</script>

{#if !friend.accepted}
	<div class="flex items-center gap-x-3 p-1.5 rounded-xl border border-transparent">
		<FriendButtonContent {...friend} sender={friend.friendship_sender_id === userStore.user?.id} />

		{#if friend.friendship_sender_id !== userStore.user?.id}
			<FriendshipControls
				friendshipID={friend.friendship_id}
				receiverID={userStore.user!.id}
				senderID={friend.friendship_sender_id}
			/>
		{/if}
	</div>
{:else}
	<button
		class={[
			'relative text-left flex items-center gap-x-3  p-1.5 rounded-xl transition-colors duration-100 border',
			friend.accepted && 'active-scale-down',
			page.url.pathname.includes(friend.channel_id!)
				? 'bg-main-900 border-main-800'
				: 'hover:bg-main-900 hover:border-main-800 border-transparent'
		]}
		onclick={() => goto(`/friends/${friend.channel_id}`)}
	>
		<FriendButtonContent {...friend} sender={friend.friendship_sender_id === userStore.user?.id} />

		<ContextMenuFriend {friend} />
	</button>
{/if}
