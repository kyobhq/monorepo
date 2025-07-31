export class Core {
  serverDialog = $state(false)
  categoryDialog = $state(false)
  channelDialog = $state({ open: false, category_id: '' })
  friendsDialog = $state(false)
  settingsDialog = $state({ open: false, section: '' });
  serverSettingsDialog = $state({ open: false, server_id: '', section: '' });
  categorySettingsDialog = $state({ open: false, category_id: '', section: '' });
  channelSettingsDialog = $state({ open: false, channel_id: '', section: '' });
  destructiveDialog = $state({ open: false, title: '', subtitle: '', buttonText: '', onclick: () => { } })
}

export const coreStore = new Core();
