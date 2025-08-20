import type { User } from '$lib/types/types';
import { backend } from './backendStore.svelte';
import { userStore } from './userStore.svelte';

const MAX_SEEN_PROFILES = 5;

export class Core {
  firstLoad = $state({ sidebar: false, serverbar: false })
  serversLoaded = $state(false);
  richInputLength = $state(0);
  serverDialog = $state(false);
  categoryDialog = $state(false);
  friendsDialog = $state(false);
  channelDialog = $state({ open: false, category_id: '' });
  modDialog = $state({ open: false, action: 'kick', server_id: '', user_id: '' });
  restrictionDialog = $state({ open: false, restriction: '', title: '', reason: '' });
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
  isUsingKeyboard = $state(false);
  profile = $state<{ open: boolean; user: User | null; el: HTMLElement | null; side: string }>({
    open: false,
    user: null,
    el: null,
    side: 'top'
  });
  seenProfiles = $state<User[]>([]);

  handleShiftUp = (e: KeyboardEvent) => {
    if (e.key === 'Shift') {
      this.pressingShift = false;
    }
  };

  handleKeyboardInput = (e: KeyboardEvent) => {
    if (e.key === 'Tab') {
      this.isUsingKeyboard = true;
      this.updateDataAttributes();
    }

    if (e.key === 'Shift') {
      this.pressingShift = true;
    }
  };

  handleMouseInput = () => {
    this.isUsingKeyboard = false;
    this.updateDataAttributes();
  };

  handlePointerInput = () => {
    this.isUsingKeyboard = false;
    this.updateDataAttributes();
  };

  handleTouchInput = () => {
    this.isUsingKeyboard = false;
    this.updateDataAttributes();
  };

  private updateDataAttributes() {
    if (typeof document !== 'undefined') {
      document.documentElement.setAttribute('data-keyboard-user', String(this.isUsingKeyboard));
    }
  }

  initializeKeyboardDetection() {
    if (typeof document !== 'undefined') {
      this.updateDataAttributes();

      document.addEventListener('keydown', this.handleKeyboardInput, true);
      document.addEventListener('keyup', this.handleShiftUp, true);

      document.addEventListener('mousedown', this.handleMouseInput, true);
      document.addEventListener('mouseup', this.handleMouseInput, true);

      document.addEventListener('pointerdown', this.handlePointerInput, true);
      document.addEventListener('pointerup', this.handlePointerInput, true);

      document.addEventListener('touchstart', this.handleTouchInput, true);
      document.addEventListener('touchend', this.handleTouchInput, true);
    }
  }

  cleanupKeyboardDetection() {
    if (typeof document !== 'undefined') {
      document.removeEventListener('keydown', this.handleKeyboardInput, true);
      document.removeEventListener('mousedown', this.handleMouseInput, true);
      document.removeEventListener('mouseup', this.handleMouseInput, true);
      document.removeEventListener('pointerdown', this.handlePointerInput, true);
      document.removeEventListener('pointerup', this.handlePointerInput, true);
      document.removeEventListener('touchstart', this.handleTouchInput, true);
      document.removeEventListener('touchend', this.handleTouchInput, true);
    }
  }

  openMyProfile(el: HTMLElement, side: 'top' | 'bottom' | 'left' | 'right' = 'top') {
    if (this.profile.user) {
      this.profile = {
        open: false,
        user: null,
        el: null,
        side: side
      };
    } else {
      // First set the user and element without opening
      this.profile = {
        open: false,
        user: userStore.user!,
        el: el,
        side: side
      };
      // Then open it on the next tick to trigger the animation
      setTimeout(() => {
        if (this.profile.user) {
          this.profile.open = true;
        }
      }, 10);
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
        // First set the user and element without opening
        this.profile = {
          open: false,
          user: cachedUser,
          el: el,
          side: side
        };
        // Then open it on the next tick to trigger the animation
        setTimeout(() => {
          if (this.profile.user) {
            this.profile.open = true;
          }
        }, 10);
        return;
      }

      const res = await backend.getUserProfile(userID);
      res.match(
        (user) => {
          // First set the user and element without opening
          this.profile = {
            open: false,
            user: user,
            el: el,
            side: side
          };
          // Then open it on the next tick to trigger the animation
          setTimeout(() => {
            if (this.profile.user) {
              this.profile.open = true;
            }
          }, 10);

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
