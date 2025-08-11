<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import DefaultSettingsDialog from '../DefaultSettingsDialog.svelte';
	import SideBarSettings from 'ui/SideBar/SideBarSettings.svelte';
	import UserSettingsProfile from './UserSettingsProfile.svelte';
	import UserSettingsPassword from './UserSettingsPassword.svelte';
	import UserSettingsEmail from './UserSettingsEmail.svelte';
	import Separator from '../Separator.svelte';
	import UserSettingsAvatar from './UserSettingsAvatar.svelte';
	import UserSettingsEmojis from './UserSettingsEmojis.svelte';

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
			<div>
				<p class="font-medium select-none">Account removal</p>
				<p class="text-sm text-main-500 select-none">
					This is not a soft delete! This action is irreversible.
				</p>
				<button
					class="text-left w-fit bg-red-400/30 border-[0.5px] border-red-400 px-2 py-1.5 text-red-400 hocus:bg-red-400 hocus:text-red-50 hover:cursor-pointer transition-colors duration-100 mt-4"
				>
					Delete account
				</button>
			</div>
		{:else if coreStore.userSettingsDialog.section === 'Profile'}
			<UserSettingsAvatar />
			<UserSettingsProfile />
		{:else if coreStore.userSettingsDialog.section === 'Emojis'}
			<UserSettingsEmojis />
		{/if}
	</div>
</DefaultSettingsDialog>
