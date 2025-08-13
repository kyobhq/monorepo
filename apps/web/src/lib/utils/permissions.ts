import type { Abilities } from "$lib/constants/permissions";
import { serverStore } from "stores/serverStore.svelte";

export function isOwner(serverID: string) {
  return serverStore.abilities[serverID].includes("OWNER")
}

export function hasPermissions(serverID: string, ...abilities: Abilities[]) {
  if (!serverStore.abilities[serverID]) return false

  return serverStore.abilities[serverID].some(ability => abilities.includes(ability)) ||
    serverStore.abilities[serverID].includes("OWNER") ||
    serverStore.abilities[serverID].includes("ADMINISTRATOR")
}
