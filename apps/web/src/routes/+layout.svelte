<script lang="ts">
	import { goto } from '$app/navigation';
	import '../app.css';
	import '@fontsource/host-grotesk/400.css';
	import '@fontsource/host-grotesk/500.css';
	import '@fontsource/host-grotesk/600.css';
	import '@fontsource/host-grotesk/700.css';
	import { backend } from 'stores/backendStore.svelte';
	import { userStore } from 'stores/userStore.svelte';
	import { onMount } from 'svelte';
	let { children } = $props();

	onMount(async () => {
		const identityRes = await backend.checkIdentity();
		identityRes.match(
			(user) => {
				userStore.user = user;
				goto('/servers');
			},
			(error) => goto('/signin')
		);
	});
</script>

{@render children()}
