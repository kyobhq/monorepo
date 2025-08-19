<script lang="ts">
	import { afterNavigate, goto } from '$app/navigation';
	import { page } from '$app/state';
	import { serverStore } from 'stores/serverStore.svelte';
	import BarSeparator from 'ui/BarSeparator/BarSeparator.svelte';
	import Channel from 'ui/Channel/Channel.svelte';
	import CollapsibleBox from 'ui/CollapsibleBox/CollapsibleBox.svelte';
	import ContextMenuSideBar from 'ui/ContextMenu/ContextMenuSideBar.svelte';
	import gsap from 'gsap';

	const currentServer = $derived(serverStore.getServer(page.params.server_id || '') || undefined);

	let sectionEl = $state<HTMLElement | undefined>(undefined);

	$effect(() => {
		if (!sectionEl) return;
		const boxes = sectionEl.querySelectorAll('.collapsible-wrapper');
		if (!boxes.length) return;
		gsap.from(boxes, {
			opacity: 0,
			scale: 0.95,
			stagger: 0.06,
			duration: 0.35,
			ease: 'power2.out'
		});
	});
</script>

{#if currentServer}
	<BarSeparator title={currentServer.name} />
	<section class="relative flex flex-col gap-y-2 p-2.5 h-full" bind:this={sectionEl}>
		<ContextMenuSideBar />
		{#each Object.values(currentServer.categories).sort((a, b) => a.position - b.position) as category (category.id)}
			<div class="collapsible-wrapper">
				<CollapsibleBox header={category.name} categoryId={category.id}>
					{#if category.channels}
						{#each Object.values(category.channels).sort((a, b) => a.position - b.position) as channel (channel.id)}
							<Channel
								id={channel.id}
								categoryId={category.id}
								type={channel.type}
								name={channel.name}
								onclick={() => goto(`/servers/${channel.server_id}/channels/${channel.id}`)}
								active={page.url.pathname.includes(channel.id)}
							/>
						{/each}
					{/if}
				</CollapsibleBox>
			</div>
		{/each}
	</section>
{/if}
