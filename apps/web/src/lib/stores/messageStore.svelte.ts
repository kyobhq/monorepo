import type { Message } from "$lib/types/types";

export class MessageStore {
  editMessage = $state<Message | null>(null)

  stopEditing() {
    this.editMessage = null;
  }
}

export const messageStore = new MessageStore();
