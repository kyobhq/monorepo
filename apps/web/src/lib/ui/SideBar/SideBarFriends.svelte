<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { onMount } from 'svelte';
	import BarSeparatorWithButton from 'ui/BarSeparator/BarSeparatorWithButton.svelte';
	import FriendButton from 'ui/FriendButton/FriendButton.svelte';

	onMount(() => {
		coreStore.firstLoad.sidebar = true;
	});
</script>

<BarSeparatorWithButton
	title={`Friends - ${userStore.getNbOnlineFriends()}`}
	buttonText="Add Friend"
	onclick={() => {
		coreStore.friendsDialog = true;
	}}
/>

<div class="flex flex-col gap-y-1 px-2.5 pt-2.5">
	{#each userStore.friends as friend (friend.id)}
		<FriendButton {friend} />
	{/each}
</div>
