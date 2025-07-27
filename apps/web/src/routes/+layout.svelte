<script lang="ts">
	import '../app.css';
	import SideBar from '../ui/SideBar/SideBar.svelte';
	import ServerBar from '../ui/ServerBar/ServerBar.svelte';
	let { children } = $props();
</script>

<div class="flex">
	<SideBar />
	<main class="flex-1">
		{@render children()}
	</main>
	<ServerBar />
</div>
