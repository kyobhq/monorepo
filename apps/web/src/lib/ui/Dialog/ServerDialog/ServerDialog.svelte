<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import DefaultDialog from '../DefaultDialog.svelte';
	import CreateServerForm from './CreateServerForm.svelte';
	import JoinServerForm from './JoinServerForm.svelte';

	const TABS: Record<string, { title: string; subtitle: string }> = {
		Create: {
			title: 'Create a server',
			subtitle: 'It can be either public or private.'
		},
		Join: {
			title: 'Join a server',
			subtitle: 'Just paste your invite link and off you go.'
		},
		Discover: {
			title: 'Discover servers',
			subtitle: 'This is where public servers appear!'
		}
	};

	let currentTab = $state('Create');
</script>

<DefaultDialog
	bind:state={coreStore.serverDialog}
	title={TABS[currentTab].title}
	subtitle={TABS[currentTab].subtitle}
	tabs={Object.keys(TABS)}
	bind:currentTab
>
	{#if currentTab === 'Create'}
		<CreateServerForm />
	{:else if currentTab === 'Join'}
		<JoinServerForm />
	{/if}
</DefaultDialog>
