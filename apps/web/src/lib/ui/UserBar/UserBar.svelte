<script lang="ts">
	import gsap from 'gsap';
	import { coreStore } from 'stores/coreStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { fly } from 'svelte/transition';
	import Gear from 'ui/icons/Gear.svelte';
	import Headphone from 'ui/icons/Headphone.svelte';
	import Microphone from 'ui/icons/Microphone.svelte';

	let userBarEl = $state<HTMLElement>();
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

	$effect(() => {
		if (!userBarEl) return;

		gsap.from(userBarEl, {
			opacity: 0,
			scale: 0.95,
			duration: 0.35,
			ease: 'power2.out'
		});
	});
</script>

{#if userStore.user}
	<div
		bind:this={userBarEl}
		class="absolute w-[calc(100%-1.25rem)] bottom-2.5 pr-2 left-1/2 -translate-x-1/2 box-style rounded-xl"
	>
		<div class="flex justify-between p-[2.5px]">
			<button
				bind:this={buttonEl}
				onclick={() => coreStore.openMyProfile(buttonEl!)}
				class="flex text-left items-center gap-x-2 active:bg-main-700/85 hover:bg-main-800 pr-2 border-[0.5px] border-transparent hover:border-main-600 transition-colors hover:cursor-pointer duration-150 rounded-lg z-[1]"
			>
				<figure class="relative h-12 w-12 highlight-border rounded-lg overflow-hidden">
					<img src={userStore.user.avatar} alt="" class="w-full h-full object-cover" />
				</figure>
				<div class="flex flex-col gap-y-0.5">
					<p class="leading-none font-medium truncate w-24">{userStore.user.display_name}</p>
					<p class="text-sm leading-none text-main-500">Connected</p>
				</div>
			</button>

			<div class="flex items-center gap-x-1 z-[1]">
				<button
					onclick={handleMute}
					class={['icon-button', userStore.mute ? 'icon-button-destructive' : 'icon-button-normal']}
				>
					<Microphone height={20} width={20} mute={userStore.mute} />
				</button>
				<button
					onclick={handleDeafen}
					class={[
						'icon-button',
						userStore.deafen ? 'icon-button-destructive' : 'icon-button-normal'
					]}
				>
					<Headphone height={20} width={20} deafen={userStore.deafen} />
				</button>
				<button onclick={openSettings} class="icon-button icon-button-normal">
					<Gear height={20} width={20} class="flex justify-center items-center" />
				</button>
			</div>
		</div>
	</div>
{/if}
