<script lang="ts">
	import ContextMenuSkeleton from './ContextMenuSkeleton.svelte';
	import ContextMenuItem from './ContextMenuItem.svelte';
	import type { Friend } from '$lib/types/types';
	import { backend } from 'stores/backendStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { logErr } from 'utils/print';
	import RemoveUserIcon from 'ui/icons/RemoveUserIcon.svelte';
	import BlockUserIcon from 'ui/icons/BlockUserIcon.svelte';
	import UpPhoneIcon from 'ui/icons/UpPhoneIcon.svelte';

	interface Props {
		friend: Friend;
	}

	let { friend }: Props = $props();

	async function handleRemoveFriend() {
		const res = await backend.removeFriend({
			friendship_id: friend.friendship_id,
			sender_id: friend.friendship_sender_id,
			receiver_id:
				friend.friendship_sender_id === userStore.user!.id ? friend.id! : userStore.user!.id,
			channel_id: friend.channel_id!
		});
		res.match(
			() => {
				userStore.removeFriend({ friendshipID: friend.friendship_id });
			},
			(err) => logErr(err)
		);
	}
</script>

<ContextMenuSkeleton>
	{#snippet contextMenuContent()}
		<ContextMenuItem Icon={UpPhoneIcon} onclick={() => {}} text="Call" />
		<ContextMenuItem
			Icon={RemoveUserIcon}
			onclick={handleRemoveFriend}
			text="Remove Friend"
			destructive
		/>
		<ContextMenuItem Icon={BlockUserIcon} onclick={() => {}} text="Block" destructive />
	{/snippet}
</ContextMenuSkeleton>
