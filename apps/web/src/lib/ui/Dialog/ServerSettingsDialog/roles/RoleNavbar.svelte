<script lang="ts">
	let { TABS, activeTab = $bindable() } = $props();
</script>

<nav>
	<ul class="flex gap-x-1">
		{#each TABS as tab, idx (idx)}
			<li>
				<button
					onclick={() => (activeTab = tab)}
					class={[
						'border-[0.5px] px-2 py-1 rounded-sm transition-colors duration-100',
						activeTab === tab
							? 'text-main-50 bg-main-800 border-main-600'
							: 'text-main-300 bg-main-900 border-main-700 hover:bg-main-800/75 hover:border-main-600/75'
					]}
				>
					{tab}
				</button>
			</li>
		{/each}
	</ul>
</nav>
