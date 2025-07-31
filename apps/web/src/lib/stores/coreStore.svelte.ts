export class Core {
  openServerDialog = $state(false)
  openCategoryDialog = $state(false)
  openChannelDialog = $state({ open: false, category_id: '' })
  openFriendsDialog = $state(false)
  openDestructiveDialog = $state({ open: false, title: '', subtitle: '', buttonText: '', onclick: () => { } })
}

export const coreStore = new Core();
