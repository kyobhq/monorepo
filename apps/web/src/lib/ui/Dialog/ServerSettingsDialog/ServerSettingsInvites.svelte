<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { expiresAt } from 'utils/time';

	const invites = $derived(
		coreStore.serverSettingsDialog.server_id
			? serverStore.getInvites(coreStore.serverSettingsDialog.server_id)
			: []
	);
</script>

<div class="flex-1 overflow-y-auto border-[0.5px] border-main-800 bg-main-900 mt-6">
	<table class="w-full border-collapse bg-main-900">
		<thead class="bg-main-900 sticky top-0 z-10">
			<tr>
				<th
					class="px-4 py-3 text-left font-semibold text-xs uppercase text-main-400 border-b border-main-700 select-none"
					>Inviter</th
				>
				<th
					class="px-4 py-3 text-left font-semibold text-xs uppercase text-main-400 border-b border-main-700 select-none"
					>Invite Code</th
				>
				<th
					class="px-4 py-3 text-left font-semibold text-xs uppercase text-main-400 border-b border-main-700 select-none"
					>Expires</th
				>
			</tr>
		</thead>
		<tbody>
			{#if invites && invites.length > 0}
				{#each invites as invite}
					<tr class="border-b border-main-800 transition-colors duration-150 hover:bg-main-850">
						<td class="px-4 py-3 align-middle">
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 overflow-hidden flex-shrink-0">
									{#if invite.creator?.avatar}
										<img
											src={invite.creator.avatar}
											alt={invite.creator.display_name || invite.creator.username}
											class="w-full h-full object-cover select-none"
										/>
									{:else}
										<div
											class="w-full h-full bg-accent-600 flex items-center justify-center text-white font-semibold text-sm select-none"
										>
											{(invite.creator?.display_name || invite.creator?.username || '?')
												.charAt(0)
												.toUpperCase()}
										</div>
									{/if}
								</div>
								<div class="flex flex-col gap-y-1">
									<div class="font-semibold text-white text-sm select-none leading-none">
										{invite.creator?.display_name || invite.creator?.username || 'Unknown User'}
									</div>
									{#if invite.creator?.display_name && invite.creator?.username !== invite.creator?.display_name}
										<div class="text-xs text-main-400 select-none leading-none">
											@{invite.creator.username}
										</div>
									{/if}
								</div>
							</div>
						</td>
						<td class="px-4 py-3 align-middle">
							<div class="flex items-center gap-2">
								<code class="bg-main-950 text-accent-400 px-2 py-1 text-sm font-mono select-all">
									{invite.invite_id}
								</code>
								<button
									class="text-main-400 hover:text-white transition-colors hover:cursor-pointer duration-100"
									onclick={() => navigator.clipboard.writeText(invite.invite_id)}
									title="Copy invite code"
									aria-label="Copy invite code"
								>
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
										></path>
									</svg>
								</button>
							</div>
						</td>
						<td class="px-4 py-3 align-middle">
							<span class="text-sm select-none">{expiresAt(invite.expire_at)}</span>
						</td>
					</tr>
				{/each}
			{:else}
				<tr>
					<td colspan="3" class="text-center text-main-400 py-6 italic select-none">
						No invites found
					</td>
				</tr>
			{/if}
		</tbody>
	</table>
</div>
