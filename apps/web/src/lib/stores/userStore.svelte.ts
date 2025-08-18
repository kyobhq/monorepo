import type { Channel, Emoji, Friend, User } from '$lib/types/types';

export class UserStore {
  user = $state<User>();
  friends = $state<Friend[]>([]);
  emojis = $state<Emoji[]>([]);
  pinned_channels = $state<Channel[]>([]);
  mute = $state(false);
  deafen = $state(false);
  setupComplete = $state(false);

  setFriendStatus(friendID: string, status: string) {
    const friend = this.friends.find((friend) => friend.id === friendID)
    if (friend) friend.status = status
  }

  getNbOnlineFriends() {
    return this.friends.filter(friend => friend.status !== "offline" && friend.accepted).length
  }

  acceptFriend(friendshipID: string, channelID: string) {
    const friendship = this.friends.find((friend) => friend.friendship_id === friendshipID)
    if (friendship) {
      friendship.accepted = true
      friendship.channel_id = channelID
    }
  }

  removeFriend({ friendshipID, userID }: { friendshipID?: string, userID?: string }) {
    if (friendshipID) {
      this.friends = this.friends.filter(friend => friend.friendship_id !== friendshipID)
    } else if (userID) {
      this.friends = this.friends.filter(friend => friend.id !== userID)
    }
  }
}

export const userStore = new UserStore();
