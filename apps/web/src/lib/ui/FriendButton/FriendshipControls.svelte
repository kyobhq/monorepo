<script lang="ts">
	import Check from 'ui/icons/Check.svelte';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';

	function acceptFriend(e: MouseEvent) {
		e.stopImmediatePropagation();
	}

	function refuseFriend(e: MouseEvent) {
		e.stopImmediatePropagation();
	}
</script>

<div class="flex gap-x-1 ml-auto">
	<button
		class="bg-green-400/20 border-[0.5px] border-green-400 text-green-400 flex justify-center items-center p-1 hover:bg-green-400/30 active:bg-green-400/40 duration-100 transition-colors rounded-sm"
		aria-label="Accept Friend"
	>
		<Check height={16} width={16} />
	</button>

	<button
		class="bg-red-400/20 border-[0.5px] border-red-400 text-red-400 flex justify-center items-center p-1 hover:bg-red-400/35 active:bg-red-400/40 duration-100 transition-colors rounded-sm"
		aria-label="Refuse Friend"
	>
		<PlusIcon height={16} width={16} class="rotate-45" />
	</button>
</div>
