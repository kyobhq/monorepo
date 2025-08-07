import type { AbilitiesType } from "$lib/types/types";
import { serverStore } from "stores/serverStore.svelte";
import { userStore } from "stores/userStore.svelte";

export function checkActionPermission(serverID: string, action: AbilitiesType) {
  const server = serverStore.getServer(serverID)
  if (server.owner_id === userStore.user?.id) return true
}
