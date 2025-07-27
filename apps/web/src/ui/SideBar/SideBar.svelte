<script lang="ts">
	import { goto } from '$app/navigation';
	import ServerButton from '../ServerButton/ServerButton.svelte';
	import Input from '@kyob/design-system/input';
	import IconsTabBar from '../IconsTabBar/IconsTabBar.svelte';
	import MagnifyingGlass from '../icons/MagnifyingGlass.svelte';
	import BarSeparator from '../BarSeparator/BarSeparator.svelte';
	import CollapsibleBox from '../CollapsibleBox/CollapsibleBox.svelte';
	import Channel from '../Channel/Channel.svelte';
	import { page } from '$app/state';
	import { TABS, SERVERS, PINNED_CHANNELS } from '../../constants/constants';
</script>

<aside
	class="bg-main-975 min-h-screen border-r-[0.5px] border-r-main-700 w-[19.5rem] overflow-hidden"
>
	<section class="flex flex-col gap-y-3 mb-2.5">
		<div class="flex gap-x-2 px-2.5 pt-2.5">
			<Input Icon={MagnifyingGlass} placeholder="Search" />
			<IconsTabBar
				tabs={TABS}
				onclick={(href: string) => goto(href)}
				activeTab={page.url.pathname}
			/>
		</div>
		<div
			class="flex gap-x-2.5 pl-2.5 shrink-0 relative after:absolute after:top-0 after:right-0 after:w-16 after:h-full after:bg-gradient-to-l after:from-main-975 after:to-transparent after:pointer-events-none"
		>
			{#each SERVERS as server (server.id)}
				<ServerButton
					image={server.image}
					onclick={() => goto(server.href)}
					active={page.url.pathname.includes(server.href)}
				/>
			{/each}
		</div>
		<div class="px-2.5">
			<CollapsibleBox header="Pinned channels">
				{#each PINNED_CHANNELS as channel (channel.id)}
					{@const channelHref = `/servers/${channel.server_id}/channels/${channel.id}`}

					<Channel
						type={channel.type}
						title={channel.title}
						subtitle={channel.subtitle}
						onclick={() => goto(channelHref)}
						active={page.url.pathname.includes(channelHref)}
					/>
				{/each}
			</CollapsibleBox>
		</div>
	</section>
	<BarSeparator title="The Valley" />
	<section class="flex flex-col gap-y-2 p-2.5">
		<CollapsibleBox header="general">
			<Channel type="textual" title="General" onclick={() => {}} />
		</CollapsibleBox>
		<CollapsibleBox header="cool stuff">
			<Channel type="textual" title="General" onclick={() => {}} />
		</CollapsibleBox>
		<CollapsibleBox header="vocals">
			<Channel type="voice" title="General" onclick={() => {}} />
			<Channel type="voice" title="Cowork" onclick={() => {}} />
		</CollapsibleBox>
	</section>
</aside>
