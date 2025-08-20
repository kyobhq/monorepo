<script lang="ts">
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import { joinedAt } from 'utils/time';

	const members = $derived(
		coreStore.serverSettingsDialog.server_id
			? serverStore.getMembers(coreStore.serverSettingsDialog.server_id)
			: []
	);
</script>

<div class="flex-1 overflow-y-auto border-[0.5px] border-main-800 bg-main-900 mt-6 rounded-md">
	<table class="w-full border-collapse bg-main-900">
		<thead class="bg-main-900 sticky top-0 z-10">
			<tr>
				<th
					class="px-4 py-3 text-left font-semibold text-xs uppercase text-main-400 border-b border-main-700 select-none"
					>Name</th
				>
				<th
					class="px-4 py-3 text-left font-semibold text-xs uppercase text-main-400 border-b border-main-700 select-none"
					>Member Since</th
				>
				<th
					class="px-4 py-3 text-left font-semibold text-xs uppercase text-main-400 border-b border-main-700 select-none"
					>Joined Kyob</th
				>
				<th
					class="px-4 py-3 text-left font-semibold text-xs uppercase text-main-400 border-b border-main-700 select-none"
					>Roles</th
				>
			</tr>
		</thead>
		<tbody>
			{#if members && members.length > 0}
				{#each members as member (member.id)}
					<tr class="border-b border-main-800 transition-colors duration-150 hover:bg-main-850">
						<td class="px-4 py-3 align-middle">
							<div class="flex items-center gap-3">
								<div class="w-10 h-10 overflow-hidden flex-shrink-0">
									{#if member.avatar}
										<img
											src={member.avatar}
											alt={member.display_name || member.username}
											class="w-full h-full object-cover select-none"
										/>
									{:else}
										<div
											class="w-full h-full bg-accent-600 flex items-center justify-center text-white font-semibold text-sm select-none"
										>
											{(member.display_name || member.username || '?').charAt(0).toUpperCase()}
										</div>
									{/if}
								</div>
								<div class="flex flex-col gap-y-1">
									<div class="font-semibold text-white text-sm select-none leading-none">
										{member.display_name || member.username}
									</div>
									{#if member.display_name && member.username !== member.display_name}
										<div class="text-xs text-main-400 select-none leading-none">
											@{member.username}
										</div>
									{/if}
								</div>
							</div>
						</td>
						<td class="px-4 py-3 align-middle">
							<span class="text-sm select-none">{joinedAt(member.joined_server)}</span>
						</td>
						<td class="px-4 py-3 align-middle">
							<span class="text-sm select-none">{joinedAt(member.joined_kyob)}</span>
						</td>
						<td class="px-4 py-3 align-middle">
							{#if member.roles && member.roles.length > 0}
								{#each member.roles as roleID (roleID)}
									{@const role = serverStore.getRole(
										coreStore.serverSettingsDialog.server_id,
										roleID
									)}

									{#if role && role.name !== 'Default Permissions'}
										<div
											class="py-1 px-2 rounded-md w-fit"
											style="background-color: {role.color}26; color: {role.color}"
										>
											{role.name}
										</div>
									{/if}
								{/each}
							{:else}
								<p class="text-main-400 text-sm select-none">No roles</p>
							{/if}
						</td>
					</tr>
				{/each}
			{:else}
				<tr>
					<td colspan="4" class="text-center text-main-400 py-6 italic select-none">
						No members found
					</td>
				</tr>
			{/if}
		</tbody>
	</table>
</div>
