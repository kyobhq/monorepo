<script lang="ts">
	import type { User } from '$lib/types/types';

	interface Props {
		user: User;
	}

	let { user }: Props = $props();
</script>

<figure class="relative z-[4] h-16 w-16 highlight-border rounded-lg overflow-hidden">
	<img
		src={user.avatar}
		alt="{user.username}'s avatar"
		class="h-full w-full transform-gpu object-cover"
	/>
</figure>

<div class="flex gap-x-1 items-center mt-2.5 justify-between">
	<div class="flex flex-col gap-y-0.5">
		<p class="relative z-[4] leading-none text-2xl font-bold text-white/75 mix-blend-plus-lighter">
			{user.display_name}
		</p>
		<p class="relative z-[4] leading-none text-sm text-white/50 mix-blend-plus-lighter w-fit">
			@{user.username}
		</p>
	</div>
</div>
