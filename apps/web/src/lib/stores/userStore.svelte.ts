import type { Channel, Emoji, Friend, User } from '$lib/types/types';

export class UserStore {
  user = $state<User>();
  friends = $state<Friend[]>([]);
  emojis = $state<Emoji[]>([]);
  pinned_channels = $state<Channel[]>([]);
  mute = $state(false);
  deafen = $state(false);
  setupComplete = $state(false);

  acceptFriend(friendshipID: string, channelID: string) {
    const friendship = this.friends.find((friend) => friend.friendship_id === friendshipID)
    if (friendship) {
      friendship.accepted = true
      friendship.channel_id = channelID
    }
  }

  removeFriend(friendshipID: string) {
    this.friends = this.friends.filter(friend => friend.friendship_id !== friendshipID)
  }
}

export const userStore = new UserStore();
