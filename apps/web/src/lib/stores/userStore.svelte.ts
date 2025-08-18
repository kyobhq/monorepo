import type { Channel, Emoji, Friend, User } from '$lib/types/types';

export class UserStore {
  user = $state<User>();
  friends = $state<Friend[]>([]);
  emojis = $state<Emoji[]>([]);
  pinned_channels = $state<Channel[]>([]);
  mute = $state(false);
  deafen = $state(false);
  setupComplete = $state(false);
}

export const userStore = new UserStore();
