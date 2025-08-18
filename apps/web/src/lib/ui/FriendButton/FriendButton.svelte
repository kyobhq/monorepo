<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Friend } from '$lib/types/types';
	import { userStore } from 'stores/userStore.svelte';
	import FriendButtonContent from './FriendButtonContent.svelte';
	import FriendshipControls from './FriendshipControls.svelte';

	interface Props {
		friend: Friend;
	}

	let { friend }: Props = $props();
</script>

{#if !friend.accepted}
	<div class="flex items-center gap-x-3 p-1.5 rounded-xl border border-transparent">
		<FriendButtonContent {...friend} sender={friend.friendship_sender_id === userStore.user?.id} />

		{#if friend.friendship_sender_id !== userStore.user?.id}
			<FriendshipControls />
		{/if}
	</div>
{:else}
	<button
		class={[
			'text-left flex items-center gap-x-3 hover:bg-main-900 p-1.5 rounded-xl transition-colors duration-100 border border-transparent hover:border-main-800',
			friend.accepted && 'active-scale-down'
		]}
		onclick={() => goto(`/friends/${friend.channel_id}`)}
	>
		<FriendButtonContent {...friend} sender={friend.friendship_sender_id === userStore.user?.id} />
	</button>
{/if}
