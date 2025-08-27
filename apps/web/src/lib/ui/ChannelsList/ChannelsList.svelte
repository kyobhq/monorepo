<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import type { Category, Server } from '$lib/types/types';
	import { coreStore } from 'stores/coreStore.svelte';
	import Channel from 'ui/Channel/Channel.svelte';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
	import ContextMenuSideBar from 'ui/ContextMenu/ContextMenuSideBar.svelte';
	import gsap from 'gsap';
	import { serverStore } from 'stores/serverStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { channelStore } from 'stores/channelStore.svelte';
	import { keyByProperty } from '$lib/utils/arrays';

	interface Props {
		server: Server;
	}

	let { server }: Props = $props();
	let sectionEl = $state<HTMLElement | undefined>(undefined);

	const accessibleCategories = $derived.by(() => {
		const userID = userStore.user?.id;
		if (!userID) return [];

		const categories = serverStore.getServer(server.id).categories;
		const result: Category[] = [];

		for (const category of Object.values(categories)) {
			const channels = channelStore.getCategoryChannels(server.id, category.id);
			const accessibleChannels = channels.filter((channel) =>
				channelStore.hasChannelAccess(channel, userID)
			);

			if (accessibleChannels.length > 0) {
				result.push({
					...category,
					channels: keyByProperty(accessibleChannels, 'id')
				});
			}
		}

		return result;
	});

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
	class="relative flex flex-col gap-y-2 px-1.5 pt-1.5 pb-[9rem] h-full overflow-auto"
	bind:this={sectionEl}
>
	<ContextMenuSideBar />
	{#each accessibleCategories.sort((a, b) => a.position - b.position) as category (category.id)}
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
