<script lang="ts">
	import { Popover } from 'bits-ui';
	import { coreStore } from 'stores/coreStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import CustomPopoverContent from 'ui/Popover/CustomPopoverContent.svelte';

	$effect(() => {
		if (!coreStore.profile.user) return;
		const userColor = coreStore.profile.user.main_color
			? coreStore.profile.user.main_color
			: '#121214';

		document.documentElement.style.setProperty('--user-color-75', `rgba(${userColor}, 0.75)`);
		document.documentElement.style.setProperty('--user-color-80', `rgba(${userColor}, 0.80)`);
		document.documentElement.style.setProperty('--user-color-85', `rgba(${userColor}, 0.85)`);
		document.documentElement.style.setProperty('--user-color-90', `rgba(${userColor}, 0.90)`);
		document.documentElement.style.setProperty('--user-color-95', `rgba(${userColor}, 0.95)`);
		document.documentElement.style.setProperty('--user-color-98', `rgba(${userColor}, 0.98)`);
		document.documentElement.style.setProperty('--user-color', `rgba(${userColor}, 1)`);
	});
</script>

<Popover.Root
	open={coreStore.profile.open}
	onOpenChange={(s) => {
		if (!s) {
			setTimeout(() => {
				coreStore.closeProfile();
			}, 200);
		}
	}}
>
	<Popover.Trigger class="absolute" />
	<CustomPopoverContent
		class="relative z-[999] w-[20rem] p-0"
		align="start"
		side="top"
		sideOffset={4}
		y={5}
		customAnchor={coreStore.profile.el}
	>
		<div
			class="relative z-[2] h-full overflow-hidden bg-main-950 select-none pt-[6rem] pb-5 px-3.5"
		>
			<figure class="absolute z-[1] top-0 left-0 h-[15rem] w-full">
				<img
					src={coreStore.profile.user!.banner}
					alt="{coreStore.profile.user!.username}'s banner"
					class="h-full w-full transform-gpu object-cover"
				/>
				<div class="user-profile-gradient"></div>
			</figure>

			<figure class="relative z-[4] h-16 w-16 highlight-border">
				<img
					src={coreStore.profile.user!.avatar}
					alt="{coreStore.profile.user!.username}'s avatar"
					class="h-full w-full transform-gpu object-cover"
				/>
			</figure>

			<div class="flex gap-x-1"></div>

			<div class="flex gap-x-1 items-center mt-2.5 justify-between">
				<div class="flex flex-col gap-y-0.5">
					<p
						class="relative z-[4] leading-none text-2xl font-bold text-white/75 mix-blend-plus-lighter"
					>
						{coreStore.profile.user!.display_name}
					</p>
					<p class="relative z-[4] leading-none text-sm text-white/50 mix-blend-plus-lighter w-fit">
						@{coreStore.profile.user!.username}
					</p>
				</div>

				{#if coreStore.profile.user?.id !== userStore.user?.id}
					<button
						class="relative z-[4] mix-blend-plus-lighter bg-white/20 border border-white/30 px-2 py-1 text-sm text-white/80 font-medium"
					>
						Add friend
					</button>
				{/if}
			</div>
		</div>
	</CustomPopoverContent>
</Popover.Root>
