<script lang="ts">
	import Input from 'ui/Input/Input.svelte';
	import BarSeparator from 'ui/BarSeparator/BarSeparator.svelte';
	import MagnifyingGlass from 'ui/icons/MagnifyingGlass.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { page } from '$app/state';
	import MembersList from 'ui/MembersList/MembersList.svelte';

	const currentServer = $derived(serverStore.getServer(page.params.server_id!));
</script>

{#if currentServer}
	<div class="bg-main-975 border-l-[0.5px] border-l-main-800 w-[16rem] overflow-hidden h-screen">
		<section class="p-2.5">
			<Input Icon={MagnifyingGlass} placeholder={`Search ${currentServer.name}`} />
		</section>
		<BarSeparator title={`Members - ${serverStore.memberCount}`} />
		<MembersList />
	</div>
{/if}
