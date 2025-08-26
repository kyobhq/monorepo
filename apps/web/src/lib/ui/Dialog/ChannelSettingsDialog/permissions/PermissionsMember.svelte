<script lang="ts">
	import Check from 'ui/icons/Check.svelte';

	let { id, avatar, display_name } = $props();
</script>

<li class="flex items-center justify-between hover:bg-main-800 pl-1 pr-3 py-1 rounded-xl">
	<div class="flex items-center gap-x-2">
		<figure class="h-8 w-8 rounded-lg overflow-hidden">
			<img src={avatar} alt="" />
		</figure>
		<p class="font-medium">
			{display_name}
		</p>
	</div>

	<input
		checked={false}
		id="activate-member-{id}"
		type="checkbox"
		class="appearance-none h-0 w-0 opacity-0 peer"
	/>
	<label
		for="activate-member-{id}"
		class="h-5 w-5 bg-main-800 border-[0.5px] border-main-600 hover:cursor-pointer peer-checked:border-accent-lighter peer-checked:bg-accent transition-colors duration-75 flex justify-center items-center text-main-800 peer-checked:text-main-50 rounded-md"
	>
		<Check height={14} width={14} />
	</label>
</li>
