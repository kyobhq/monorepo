<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import ContextMenuUser from 'ui/ContextMenu/ContextMenuUser.svelte';
	import AnimatedAvatar from 'ui/AnimatedAvatar/AnimatedAvatar.svelte';

	let { author, hoverAvatar = $bindable() } = $props();
	let avatarEl = $state<HTMLButtonElement>();
</script>

<button
	bind:this={avatarEl}
	class="h-12 w-12 relative highlight-border z-[1] mb-1 select-none shrink-0 hover:after:border-main-50/75 active:after:border-main-50/50 hover:cursor-pointer rounded-xl overflow-hidden"
	onclick={() => {
		if (author.id === userStore.user!.id) {
			coreStore.openMyProfile(avatarEl!, 'right');
		} else {
			coreStore.openProfile(author.id!, avatarEl!, 'right');
		}
	}}
>
	<AnimatedAvatar src={author.avatar} alt="" class="w-full h-full" hover={hoverAvatar} />
	<ContextMenuUser memberID={author.id!} />
</button>
