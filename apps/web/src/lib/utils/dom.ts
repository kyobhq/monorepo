export function searchValidMessageParent(target: HTMLElement) {
  let parent: HTMLElement = target;

  while (parent.parentElement) {
    if (parent.id.includes('message') && parent.dataset.authorId) {
      return { id: parent.id, author: parent.dataset.authorId };
    }
    parent = parent.parentElement;
  }

  return { id: null, author: null };
}
