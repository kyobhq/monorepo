<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import SideBarSettings from 'ui/SideBar/SideBarSettings.svelte';
	import DefaultSettingsDialog from '../DefaultSettingsDialog/DefaultSettingsDialog.svelte';
	import ServerSettingsAvatar from './avatar/ServerSettingsAvatar.svelte';
	import ServerSettingsInvites from './invites/ServerSettingsInvites.svelte';
	import ServerSettingsMembers from './members/ServerSettingsMembers.svelte';
	import ServerSettingsProfile from './profile/ServerSettingsProfile.svelte';
	import ServerSettingsRoles from './roles/ServerSettingsRoles.svelte';
	import { hasPermissions } from 'utils/permissions';

	const serverID = $derived(coreStore.serverSettingsDialog.server_id);
	const editingServer = $derived(serverStore.getServer(serverID));
	let initialized = $state(false);
	let container = $state<HTMLDivElement>();
	let currentSection = $derived(coreStore.serverSettingsDialog.section);

	$effect(() => {
		if (container && coreStore.serverSettingsDialog.section) {
			container.scrollTo({ top: 0 });
		}
	});
</script>

<DefaultSettingsDialog bind:state={coreStore.serverSettingsDialog.open} bind:initialized>
	<SideBarSettings
		settings={['Server Profile', 'Members', 'Roles', 'Invites', 'Bans']}
		navigationFn={(setting) => (coreStore.serverSettingsDialog.section = setting)}
		activeSection={currentSection}
	/>
	<div
		bind:this={container}
		class={[
			'flex flex-col w-full h-full overflow-auto',
			currentSection !== 'Roles' ? 'px-8 pt-6 pb-16' : ''
		]}
	>
		{#if currentSection === 'Server Profile' && hasPermissions(serverID, 'MANAGE_SERVER')}
			<h3 class="text-xl font-semibold select-none">{currentSection}</h3>
			<ServerSettingsAvatar server={editingServer} />
			<ServerSettingsProfile server={editingServer} />
		{:else if currentSection === 'Members'}
			<h3 class="text-xl font-semibold select-none">{currentSection}</h3>
			<ServerSettingsMembers />
		{:else if currentSection === 'Roles' && hasPermissions(serverID, 'MANAGE_ROLES')}
			<ServerSettingsRoles />
		{:else if currentSection === 'Invites'}
			<h3 class="text-xl font-semibold select-none">{currentSection}</h3>
			<ServerSettingsInvites />
		{:else if currentSection === 'Bans' && hasPermissions(serverID, 'BAN_MEMBERS')}
			<h3 class="text-xl font-semibold select-none">{currentSection}</h3>
			Bans server
		{:else}
			No permissions
		{/if}
	</div>
</DefaultSettingsDialog>
