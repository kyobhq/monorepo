<script lang="ts">
	import type { Role } from '$lib/types/types';
	import RoleEdition from './RoleEdition.svelte';
	import RolesSidebar from './RolesSidebar.svelte';

	let selectedRole = $state<Role | null>(null);
	let initialized = $state(false);
</script>

<div class="flex h-full">
	<RolesSidebar bind:selectedRole bind:initialized />
	{#if selectedRole}
		<RoleEdition bind:selectedRole bind:initialized />
	{:else}
		No role selected
	{/if}
</div>
