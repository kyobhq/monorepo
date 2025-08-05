import type { User } from "$lib/types/types";
import { userStore } from "./userStore.svelte";

export class Core {
  serverDialog = $state(false)
  categoryDialog = $state(false)
  channelDialog = $state({ open: false, category_id: '' })
  friendsDialog = $state(false)
  settingsDialog = $state({ open: false, section: '' });
  serverSettingsDialog = $state({ open: false, server_id: '', section: '' });
  categorySettingsDialog = $state({ open: false, category_id: '', section: '' });
  channelSettingsDialog = $state({ open: false, channel_id: '', section: '' });
  userSettingsDialog = $state({ open: false, section: '' });
  destructiveDialog = $state({ open: false, title: '', subtitle: '', buttonText: '', onclick: () => { } })
  pressingShift = $state(false)
  profile = $state<{ open: boolean, user: User | null, el: HTMLElement | null }>({ open: false, user: null, el: null });

  handleShiftDown = (e: KeyboardEvent) => {
    if (e.key === 'Shift') {
      this.pressingShift = true
    }
  }

  handleShiftUp = (e: KeyboardEvent) => {
    if (e.key === 'Shift') {
      this.pressingShift = false
    }
  }

  openMyProfile(el: HTMLElement) {
    if (this.profile.user) {
      this.profile = {
        open: false,
        user: null,
        el: null
      }
    } else {
      this.profile = {
        open: true,
        user: userStore.user!,
        el: el
      }
    }
  }

  closeProfile() {
    this.profile = {
      open: false,
      user: null,
      el: null
    }
  }
}

export const coreStore = new Core();
