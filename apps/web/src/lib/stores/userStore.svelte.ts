import type { Channel, User } from "$lib/types/types";

export class UserStore {
  user = $state<User>();
  emojis = $state([])
  pinned_channels = $state<Channel[]>([])
  mute = $state(false);
  deafen = $state(false);

}

export const userStore = new UserStore()
