<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import type { Server } from '$lib/types/types';
	import { coreStore } from 'stores/coreStore.svelte';
	import Channel from 'ui/Channel/Channel.svelte';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
	import ContextMenuSideBar from 'ui/ContextMenu/ContextMenuSideBar.svelte';
	import gsap from 'gsap';

	interface Props {
		server: Server;
	}

	let { server }: Props = $props();
	let sectionEl = $state<HTMLElement | undefined>(undefined);

	$effect(() => {
		if (!sectionEl || coreStore.firstLoad.sidebar) return;
		const boxes = sectionEl.querySelectorAll('.collapsible-wrapper');
		if (!boxes.length) return;
		gsap.from(boxes, {
			opacity: 0,
			y: 8,
			stagger: 0.06,
			duration: 0.35,
			ease: 'power2.out',
			onComplete: () => {
				coreStore.firstLoad.sidebar = true;
			}
		});
	});
</script>

<section
	class="relative flex flex-col gap-y-2 px-2.5 pt-2.5 pb-[9rem] h-full overflow-auto"
	bind:this={sectionEl}
>
	<ContextMenuSideBar />
	{#each Object.values(server.categories).sort((a, b) => a.position - b.position) as category (category.id)}
		<div class="collapsible-wrapper">
			<CollapsibleBox header={category.name} categoryId={category.id}>
				{#if category.channels}
					{#each Object.values(category.channels).sort((a, b) => a.position - b.position) as channel (channel.id)}
						{@const unread = channel.last_message_sent !== channel.last_message_read}
						<Channel
							id={channel.id}
							categoryId={category.id}
							type={channel.type}
							name={channel.name}
							onclick={() => goto(`/servers/${channel.server_id}/channels/${channel.id}`)}
							active={page.url.pathname.includes(channel.id)}
							{unread}
							mentions={channel.last_mentions?.length || 0}
						/>
					{/each}
				{/if}
			</CollapsibleBox>
		</div>
	{/each}
</section>
