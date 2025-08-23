import type { SuggestionProps } from '@tiptap/suggestion';

export class Editor {
	currentChannel = $state('');
	currentInput = $state<'main' | 'edit'>('main');
	listOpen = $state(false);
	mentionProps = $state<SuggestionProps | null>();
	mentionsListEl = $state<any>(null);
	emojiProps = $state<SuggestionProps | null>();
	emojisListEl = $state<any>();
	menuHeight = $state(0);
}

export const editorStore = new Editor();
