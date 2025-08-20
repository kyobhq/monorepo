<script lang="ts">
	import ContextMenuSkeleton from './ContextMenuSkeleton.svelte';
	import ContextMenuItem from './ContextMenuItem.svelte';
	import type { Member } from '$lib/types/types';
	import { userStore } from 'stores/userStore.svelte';
	import { hasPermissions } from 'utils/permissions';
	import { page } from '$app/state';
	import { coreStore } from 'stores/coreStore.svelte';
	import UpPhoneIcon from 'ui/icons/UpPhoneIcon.svelte';
	import EditIcon from 'ui/icons/EditIcon.svelte';
	import BlockUserIcon from 'ui/icons/BlockUserIcon.svelte';
	import HammerIcon from 'ui/icons/HammerIcon.svelte';
	import CrossUserIcon from 'ui/icons/CrossUserIcon.svelte';
	import MessageIcon from 'ui/icons/MessageIcon.svelte';
	import { serverStore } from 'stores/serverStore.svelte';

	interface Props {
		memberID: string;
	}

	let { memberID }: Props = $props();

	const serverID = $derived(page.params.server_id || '');
	const member = $derived(serverStore.getMember(serverID, memberID)!);
	let isFriend = $derived(userStore.friends.find((friend) => friend.id === member.id));

	async function handleBan() {
		if (!page.params.server_id || !member.id) return;
		coreStore.modDialog = {
			open: true,
			action: 'ban',
			server_id: page.params.server_id,
			user_id: member.id
		};
	}

	async function handleKick() {
		if (!page.params.server_id || !member.id) return;
		coreStore.modDialog = {
			open: true,
			action: 'kick',
			server_id: page.params.server_id,
			user_id: member.id
		};
	}
</script>

<ContextMenuSkeleton>
	{#snippet contextMenuContent()}
		{#if isFriend}
			<ContextMenuItem Icon={UpPhoneIcon} onclick={() => {}} text="Call" />
			<ContextMenuItem Icon={MessageIcon} onclick={() => {}} text="Message" />
		{/if}
		{#if member.id === userStore.user?.id || hasPermissions(serverID, 'MANAGE_NICKNAMES')}
			<ContextMenuItem Icon={EditIcon} onclick={() => {}} text="Change Nickname" />
		{/if}
		{#if member.id !== userStore.user?.id}
			<ContextMenuItem Icon={BlockUserIcon} onclick={() => {}} text="Block" destructive />

			{#if hasPermissions(serverID, 'KICK_MEMBERS')}
				<ContextMenuItem Icon={CrossUserIcon} onclick={handleKick} text="Kick" destructive />
			{/if}

			{#if hasPermissions(serverID, 'BAN_MEMBERS')}
				<ContextMenuItem Icon={HammerIcon} onclick={handleBan} text="Ban" destructive />
			{/if}
		{/if}
	{/snippet}
</ContextMenuSkeleton>
