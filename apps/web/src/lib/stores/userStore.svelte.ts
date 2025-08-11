import type { Channel, Emoji, User } from '$lib/types/types';

export class UserStore {
  user = $state<User>();
  emojis = $state<Emoji[]>([]);
  pinned_channels = $state<Channel[]>([]);
  mute = $state(false);
  deafen = $state(false);
  setupComplete = $state(false);
}

export const userStore = new UserStore();
