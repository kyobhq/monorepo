<script lang="ts">
	import '../app.css';
	import '@fontsource/host-grotesk/400.css';
	import '@fontsource/host-grotesk/500.css';
	import '@fontsource/host-grotesk/600.css';

	import SideBar from 'ui/SideBar/SideBar.svelte';
	import ServerBar from 'ui/ServerBar/ServerBar.svelte';
	let { children } = $props();
</script>

<div class="flex">
	<SideBar />
	<main class="flex-1 relative">
		{@render children()}
	</main>
	<ServerBar />
</div>
