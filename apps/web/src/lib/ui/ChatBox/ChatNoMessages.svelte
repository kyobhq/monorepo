<script lang="ts">
	import StraightFaceEmoji from 'ui/icons/StraightFaceEmoji.svelte';

	let { channelName } = $props();
</script>

<div class="h-full w-full flex items-center justify-center text-2xl flex-col gap-y-4 text-main-700">
	<StraightFaceEmoji height={128} width={128} />
	<p>No messages in <span class="font-medium">{channelName}</span> yet</p>
</div>
