<script lang="ts">
	import type { Role } from '$lib/types/types';
	import { backend } from 'stores/backendStore.svelte';
	import { coreStore } from 'stores/coreStore.svelte';
	import { serverStore } from 'stores/serverStore.svelte';
	import PlusIcon from 'ui/icons/PlusIcon.svelte';
	import { hasPermissions } from 'utils/permissions';
	import { logErr } from 'utils/print';

	let { selectedRole = $bindable(), initialized = $bindable() } = $props();
	let serverID = $derived(coreStore.serverSettingsDialog.server_id);
	let roles = $derived(serverStore.getRoles(serverID));
	let userTopRolePosition = $derived.by(() => {
		let topRolePosition: number = 9999;
		const userRoles = serverStore.getUserRoles(serverID);
		for (const userRole of userRoles) {
			const role = serverStore.getRole(serverID, userRole);

			if (role && role.position < topRolePosition) {
				topRolePosition = role?.position;
			}
		}

		return topRolePosition;
	});

	let draggedRole: any = $state(null);
	let draggedIndex = $state(-1);
	let dragOverIndex = $state(-1);

	function handleDragStart(event: DragEvent, role: any, index: number) {
		draggedRole = role;
		draggedIndex = index;
		if (event.dataTransfer) {
			event.dataTransfer.effectAllowed = 'move';

			const buttonElement = (event.target as HTMLElement).parentElement;
			if (buttonElement) {
				event.dataTransfer.setDragImage(buttonElement, 0, 0);
			}
		}
	}

	function handleDragOver(event: DragEvent, index: number) {
		event.preventDefault();
		if (draggedRole.position === index) return;

		dragOverIndex = index;
		if (event.dataTransfer) {
			event.dataTransfer.dropEffect = 'move';
		}
	}

	function handleDragLeave() {
		dragOverIndex = -1;
	}

	function handleDrop(event: DragEvent, targetIndex: number) {
		event.preventDefault();
		if (draggedRole.position > targetIndex && !hasPermissions(serverID, 'ADMINISTRATOR')) return;

		if (draggedRole && draggedIndex !== targetIndex) {
			const fromPosition = draggedRole.position;
			const toPosition = roles[targetIndex].position;

			backend
				.moveRole(serverID, roles[targetIndex].id, draggedRole.id, fromPosition, toPosition)
				.match(
					() => {
						const roles = serverStore.getRoles(serverID);
						const movedRole = serverStore.getRole(serverID, draggedRole.id);
						const targetRole = serverStore.getRole(serverID, roles[targetIndex].id);

						if (movedRole && targetRole) {
							movedRole.position = targetIndex;
							targetRole.position = draggedIndex;
						}

						roles.sort((a, b) => a.position - b.position);
					},
					(err) => logErr(err)
				);
		}

		draggedRole = null;
		draggedIndex = -1;
		dragOverIndex = -1;
	}

	function handleDragEnd() {
		draggedRole = null;
		draggedIndex = -1;
		dragOverIndex = -1;
	}

	function checkEditPermission(roleToEdit: Role) {
		return (
			userTopRolePosition <= roleToEdit.position ||
			serverStore.abilities[serverID].includes('ADMINISTRATOR') ||
			serverStore.abilities[serverID].includes('OWNER')
		);
	}
</script>

<div
	class="bg-main-950 h-full border-r-[0.5px] border-r-main-700 w-[13rem] overflow-hidden p-2 shrink-0 rounded-l-md flex flex-col gap-y-1"
>
	{#each roles as role, index (role.id)}
		{@const canEdit = checkEditPermission(role)}
		<button
			class={[
				'border-[0.5px] flex items-center px-3 py-2 duration-100 relative transition-all',
				dragOverIndex === index ? 'border-t-3' : '',
				!canEdit ? 'opacity-50 hover:cursor-not-allowed!' : ''
			]}
			style="background-color: {role.color}26; border-color: {role.color}80; color: {role.color};"
			onclick={() => {
				if (!canEdit) return null;

				initialized = false;
				selectedRole = role;
			}}
			ondragover={(event) => handleDragOver(event, index)}
			ondragleave={handleDragLeave}
			ondrop={(event) => handleDrop(event, index)}
		>
			<div
				draggable="true"
				role="button"
				tabindex="0"
				class="absolute inset-0 w-full h-full"
				ondragstart={(event) => handleDragStart(event, role, index)}
				ondragend={handleDragEnd}
			></div>
			<span class="relative z-10 pointer-events-none">
				{role.name}
			</span>
		</button>
	{/each}
	<button
		class="bg-main-900 border-[0.5px] border-main-800 flex items-center justify-center py-2 text-main-500 hover:text-main-50 hover:bg-main-800 hover:border-main-600 hover:cursor-pointer transition-colors duration-100"
		onclick={() => {
			initialized = false;
			serverStore.createTemplateRole(serverID);
		}}
	>
		<PlusIcon height={20} width={20} />
	</button>
</div>
