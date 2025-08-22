<script lang="ts">
	import UserLine from 'ui/UserLine/UserLine.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import type { Member, Role } from '$lib/types/types';

	interface Props {
		member: Member;
		role?: Role;
	}

	let { member, role }: Props = $props();

	let isCurrentUser = $derived(member.id === userStore.user!.id);
	let displayName = $derived(isCurrentUser ? userStore.user!.display_name : member.display_name!);
	let avatar = $derived(isCurrentUser ? userStore.user!.avatar : member.avatar!);
	let status = $derived(isCurrentUser ? 'online' : member.status);
</script>

<UserLine id={member.id!} {status} name={displayName} {avatar} hoverable color={role?.color} />
