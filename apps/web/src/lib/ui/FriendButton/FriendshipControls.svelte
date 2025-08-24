<script lang="ts">
	import { backend } from 'stores/backendStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import Check from 'ui/icons/Check.svelte';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';
	import { logErr } from 'utils/print';

	let { friendshipID, senderID, receiverID } = $props();

	async function acceptFriend() {
		const res = await backend.acceptFriendRequest({
			friendship_id: friendshipID,
			sender_id: senderID
		});
		res.match(
			() => {},
			(error) => logErr(error)
		);
	}

	async function refuseFriend() {
		const res = await backend.removeFriend({
			friendship_id: friendshipID,
			sender_id: senderID,
			receiver_id: receiverID
		});
		res.match(
			() => userStore.removeFriend({ friendshipID: friendshipID }),
			(error) => logErr(error)
		);
	}
</script>

<div class="flex gap-x-1 ml-auto">
	<button
		class="bg-green-400/20 border-[0.5px] border-green-400 text-green-400 flex justify-center items-center p-1 hover:bg-green-400/30 active:bg-green-400/40 duration-100 transition-colors rounded-sm"
		aria-label="Accept Friend"
		onclick={acceptFriend}
	>
		<Check height={16} width={16} />
	</button>

	<button
		class="bg-red-400/20 border-[0.5px] border-red-400 text-red-400 flex justify-center items-center p-1 hover:bg-red-400/35 active:bg-red-400/40 duration-100 transition-colors rounded-sm"
		aria-label="Refuse Friend"
		onclick={refuseFriend}
	>
		<PlusIcon height={16} width={16} class="rotate-45" />
	</button>
</div>
