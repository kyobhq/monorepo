<script lang="ts">
	import { Popover } from 'bits-ui';
	import { coreStore } from 'stores/coreStore.svelte';
	import CustomPopoverContent from 'ui/Popover/CustomPopoverContent.svelte';
	import UserProfileAbout from './UserProfileAbout.svelte';
	import UserProfileFacts from './UserProfileFacts.svelte';
	import UserProfileHeader from './UserProfileHeader.svelte';
	import UserProfileLinks from './UserProfileLinks.svelte';

	const user = $derived(coreStore.profile.user);
</script>

{#if user}
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
			side={coreStore.profile.side as any}
			sideOffset={8}
			alignOffset={-3}
			y={5}
			customAnchor={coreStore.profile.el}
		>
			<div
				class="relative z-[2] h-full overflow-hidden bg-main-975 select-none pt-[6rem] pb-5 px-3.5 highlight-border rounded-md"
				style="background-image: url({user.banner}); background-size: cover; background-position: center;"
			>
				<div class="absolute left-0 top-0 w-full h-full bg-main-975/85 backdrop-blur-xl"></div>
				<figure class="absolute z-[1] top-0 left-0 h-[7.5rem] w-full">
					<img
						src={user.banner}
						alt="{user.username}'s banner"
						class="h-full w-full transform-gpu object-cover"
					/>
				</figure>

				<UserProfileHeader {user} />

				<UserProfileAbout {user} />

				{#if user.links.length! > 0 || user.facts.length! > 0}
					<div
						class="relative w-full h-[0.5px] z-[4] bg-white/40 mix-blend-plus-lighter my-4.5"
					></div>

					<UserProfileLinks {user} />

					<UserProfileFacts {user} />
				{/if}
			</div>
		</CustomPopoverContent>
	</Popover.Root>
{/if}
