<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import DefaultSettingsDialog from '../DefaultSettingsDialog/DefaultSettingsDialog.svelte';
	import SideBarSettings from 'ui/SideBar/SideBarSettings.svelte';
	import UserSettingsProfile from './profile/UserSettingsProfile.svelte';
	import UserSettingsPassword from './password/UserSettingsPassword.svelte';
	import UserSettingsEmail from './email/UserSettingsEmail.svelte';
	import Separator from 'ui/Separator/Separator.svelte';
	import UserSettingsAvatar from './avatar/UserSettingsAvatar.svelte';
	import UserSettingsEmojis from './emojis/UserSettingsEmojis.svelte';
	import DestructiveBar from 'ui/DestructiveBar/DestructiveBar.svelte';
	import DeleteAccount from './deleteAccount/DeleteAccount.svelte';

	let initialized = $state(false);
	let container = $state<HTMLDivElement>();

	$effect(() => {
		if (container && coreStore.userSettingsDialog.section) {
			container.scrollTo({ top: 0 });
		}
	});
</script>

<DefaultSettingsDialog bind:state={coreStore.userSettingsDialog.open} bind:initialized>
	<SideBarSettings
		settings={['My Account', 'Profile', 'Emojis', 'Data & Privacy', 'Voice & Video']}
		navigationFn={(setting) => (coreStore.userSettingsDialog.section = setting)}
		activeSection={coreStore.userSettingsDialog.section}
	/>
	<div bind:this={container} class="flex flex-col w-full h-full px-8 pt-6 pb-16 overflow-auto">
		<h3 class="text-xl font-semibold select-none">{coreStore.userSettingsDialog.section}</h3>
		{#if coreStore.userSettingsDialog.section === 'My Account'}
			<UserSettingsEmail />
			<Separator />
			<UserSettingsPassword />
			<Separator />
			<DeleteAccount />
		{:else if coreStore.userSettingsDialog.section === 'Profile'}
			<UserSettingsAvatar />
			<UserSettingsProfile />
		{:else if coreStore.userSettingsDialog.section === 'Emojis'}
			<UserSettingsEmojis />
		{/if}
	</div>
</DefaultSettingsDialog>
