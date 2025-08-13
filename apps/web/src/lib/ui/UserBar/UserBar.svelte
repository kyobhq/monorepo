<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import Gear from 'ui/icons/Gear.svelte';
	import Headphone from 'ui/icons/Headphone.svelte';
	import Microphone from 'ui/icons/Microphone.svelte';

	let buttonEl = $state<HTMLElement>();

	function handleMute() {
		userStore.mute = !userStore.mute;
	}

	function handleDeafen() {
		userStore.deafen = !userStore.deafen;
	}

	function openSettings() {
		coreStore.userSettingsDialog = {
			open: true,
			section: 'My Account'
		};
	}

	const buttonIcon =
		'p-1 border-[0.5px] border-transparent transition-colors duration-100 hover:cursor-pointer flex items-center justify-center rounded-sm';
	const buttonIconNormal =
		'text-main-500 hover:bg-main-800 hover:border-main-600 active:bg-main-700';
	const buttonIconDestructive =
		'text-red-400 bg-red-400/20 hover:bg-red-400/30 hover:border-red-400 active:bg-red-400/50';
</script>

{#if userStore.user}
	<div
		class="absolute w-[calc(100%-1.25rem)] bg-main-900 border-[0.5px] border-main-700 bottom-2.5 py-[3.5px] pl-[3.5px] pr-2.5 left-1/2 -translate-x-1/2 flex justify-between rounded-md"
	>
		<button
			bind:this={buttonEl}
			onclick={() => coreStore.openMyProfile(buttonEl!)}
			class="flex text-left items-center gap-x-2 active:bg-main-700/65 hover:bg-main-800 pr-2.5 border-[0.5px] border-transparent hover:border-main-600 transition-colors hover:cursor-pointer duration-100 rounded-sm"
		>
			<figure class="relative h-12 w-12 highlight-border rounded-sm overflow-hidden">
				<img src={userStore.user.avatar} alt="" class="w-full h-full object-cover" />
			</figure>
			<div class="flex flex-col gap-y-0.5">
				<p class="leading-none font-medium truncate w-24">{userStore.user.display_name}</p>
				<p class="text-sm leading-none text-main-500">Connected</p>
			</div>
		</button>

		<div class="flex items-center gap-x-1">
			<button
				onclick={handleMute}
				class={[buttonIcon, userStore.mute ? buttonIconDestructive : buttonIconNormal]}
			>
				<Microphone height={20} width={20} mute={userStore.mute} />
			</button>
			<button
				onclick={handleDeafen}
				class={[buttonIcon, userStore.deafen ? buttonIconDestructive : buttonIconNormal]}
			>
				<Headphone height={20} width={20} deafen={userStore.deafen} />
			</button>
			<button onclick={openSettings} class={[buttonIcon, buttonIconNormal]}>
				<Gear height={20} width={20} class="flex justify-center items-center" />
			</button>
		</div>
	</div>
{/if}
