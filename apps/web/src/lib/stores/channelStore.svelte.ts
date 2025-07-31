import type { Channel } from "$lib/types/types";
import { categoryStore } from "./categoryStore.svelte";
import { serverStore } from "./serverStore.svelte";

class ChannelStore {
  currentChannel = $state<Channel | null>(null)

  getFirstChannel(serverID: string) {
    const server = serverStore.getServer(serverID)
    const firstCategory = Object.values(server?.categories || {})[0]
    const firstChannel = Object.values(firstCategory?.channels || {})[0]
    return firstChannel?.id || null
  }

  getChannel(serverID: string, channelID: string) {
    const server = serverStore.getServer(serverID)
    if (!server?.categories) return null;

    for (const category of Object.values(server.categories)) {
      const channel = category.channels?.[channelID];
      if (channel) return channel
    }

    return null;
  }

  getChannelsLastPositionInCategory(serverID: string, categoryID: string) {
    const category = categoryStore.getCategory(serverID, categoryID)
    return Object.values(category?.channels || {}).length
  }

  addChannel(channel: Channel) {
    const category = categoryStore.getCategory(channel.server_id, channel.category_id)
    if (!category) return null;

    category.channels[channel.id] = channel
  }

  deleteChannel(serverID: string, categoryID: string, channelID: string) {
    const category = categoryStore.getCategory(serverID, categoryID)
    if (!category) return null;

    delete category.channels[channelID]
  }
}

export const channelStore = new ChannelStore();
