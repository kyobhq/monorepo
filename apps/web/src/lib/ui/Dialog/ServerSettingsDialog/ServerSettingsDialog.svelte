<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import DefaultSettingsDialog from '../DefaultSettingsDialog.svelte';
	import SideBarSettings from 'ui/SideBar/SideBarSettings.svelte';
	import ServerSettingsProfile from './ServerSettingsProfile.svelte';
	import ServerSettingsAvatar from './ServerSettingsAvatar.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import ServerSettingsMembers from './ServerSettingsMembers.svelte';
	import ServerSettingsInvites from './ServerSettingsInvites.svelte';

	const editingServer = $derived(serverStore.getServer(coreStore.serverSettingsDialog.server_id));
	let initialized = $state(false);
	let container = $state<HTMLDivElement>();

	$effect(() => {
		if (container && coreStore.serverSettingsDialog.section) {
			container.scrollTo({ top: 0 });
		}
	});
</script>

<DefaultSettingsDialog bind:state={coreStore.serverSettingsDialog.open} bind:initialized>
	<SideBarSettings
		settings={['Server Profile', 'Members', 'Roles', 'Invites', 'Audit Log', 'Bans', 'Auto Mod']}
		navigationFn={(setting) => (coreStore.serverSettingsDialog.section = setting)}
		activeSection={coreStore.serverSettingsDialog.section}
	/>
	<div bind:this={container} class="flex flex-col w-full h-full px-8 pt-6 pb-16 overflow-auto">
		<h3 class="text-xl font-semibold select-none">{coreStore.serverSettingsDialog.section}</h3>
		{#if coreStore.serverSettingsDialog.section === 'Server Profile'}
			<ServerSettingsAvatar server={editingServer} />
			<ServerSettingsProfile server={editingServer} />
		{:else if coreStore.serverSettingsDialog.section === 'Members'}
			<ServerSettingsMembers />
		{:else if coreStore.serverSettingsDialog.section === 'Roles'}
			Roles server
		{:else if coreStore.serverSettingsDialog.section === 'Invites'}
			<ServerSettingsInvites />
		{:else if coreStore.serverSettingsDialog.section === 'Audit Log'}
			Audit Log server
		{:else if coreStore.serverSettingsDialog.section === 'Bans'}
			Bans server
		{:else if coreStore.serverSettingsDialog.section === 'Auto Mod'}
			Auto mod
		{/if}
	</div>
</DefaultSettingsDialog>
