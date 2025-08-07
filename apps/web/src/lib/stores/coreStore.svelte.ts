import type { User } from '$lib/types/types';
import { backend } from './backendStore.svelte';
import { userStore } from './userStore.svelte';

const MAX_SEEN_PROFILES = 5;

export class Core {
  serverDialog = $state(false);
  categoryDialog = $state(false);
  friendsDialog = $state(false);
  channelDialog = $state({ open: false, category_id: '' });
  settingsDialog = $state({ open: false, section: '' });
  serverSettingsDialog = $state({ open: false, server_id: '', section: '' });
  categorySettingsDialog = $state({ open: false, category_id: '', section: '' });
  channelSettingsDialog = $state({ open: false, channel_id: '', section: '' });
  userSettingsDialog = $state({ open: false, section: '' });
  destructiveDialog = $state({
    open: false,
    title: '',
    subtitle: '',
    buttonText: '',
    onclick: () => { }
  });
  pressingShift = $state(false);
  profile = $state<{ open: boolean; user: User | null; el: HTMLElement | null; side: string }>({
    open: false,
    user: null,
    el: null,
    side: 'top'
  });
  seenProfiles = $state<User[]>([]);

  handleShiftDown = (e: KeyboardEvent) => {
    if (e.key === 'Shift') {
      this.pressingShift = true;
    }
  };

  handleShiftUp = (e: KeyboardEvent) => {
    if (e.key === 'Shift') {
      this.pressingShift = false;
    }
  };

  openMyProfile(el: HTMLElement, side: 'top' | 'bottom' | 'left' | 'right' = 'top') {
    if (this.profile.user) {
      this.profile = {
        open: false,
        user: null,
        el: null,
        side: side
      };
    } else {
      this.profile = {
        open: true,
        user: userStore.user!,
        el: el,
        side: side
      };
    }
  }

  async openProfile(
    userID: string,
    el: HTMLElement,
    side: 'top' | 'bottom' | 'left' | 'right' = 'top'
  ) {
    if (this.profile.user) {
      this.profile = {
        open: false,
        user: null,
        el: null,
        side: side
      };
    } else {
      const cachedUser = this.seenProfiles.find((profile) => profile.id === userID);
      if (cachedUser) {
        this.profile = {
          open: true,
          user: cachedUser,
          el: el,
          side: side
        };
        return;
      }

      const res = await backend.getUserProfile(userID);
      res.match(
        (user) => {
          this.profile = {
            open: true,
            user: user,
            el: el,
            side: side
          };

          if (this.seenProfiles.length === MAX_SEEN_PROFILES) this.seenProfiles.shift();
          this.seenProfiles.push(user);
        },
        (error) => {
          console.error(`${error.code}: ${error.message}`);
        }
      );
    }
  }

  closeProfile() {
    this.profile = {
      open: false,
      user: null,
      el: null,
      side: 'top'
    };
  }
}

export const coreStore = new Core();
