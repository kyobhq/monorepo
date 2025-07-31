import StarterKit from '@tiptap/starter-kit';
import { EmojisSuggestion } from './extensions/emojis/emojis';
import { CustomMention } from './extensions/mentions/mentions';
import { Placeholder, type PlaceholderOptions } from '@tiptap/extensions';
import type { Content, EditorOptions, Extension, FocusPosition } from '@tiptap/core';
import type { EditorProps } from '@tiptap/pm/view';
import { editorStore } from 'stores/editorStore.svelte';

interface EditorConfigOptions {
  element: Element;
  autofocus?: FocusPosition;
  content?: Content;
  placeholder?: Extension<PlaceholderOptions, any>;
  additionalExtensions?: Extension[];
  editorProps?: EditorProps<any>;
  onTransaction?: () => void;
  onEnterPress?: () => void;
  onEscapePress?: () => void;
  onBlur?: () => void;
  onFocus?: () => void;
}

export function getMessageExtensions() {
  return [
    StarterKit.configure({
      gapcursor: false,
      dropcursor: false,
      heading: false,
      orderedList: false,
      bulletList: false,
      blockquote: false,
      link: {
        HTMLAttributes: {
          class: 'rich-link',
          role: 'button',
          rel: 'noopener noreferrer',
        },
        defaultProtocol: 'https'
      }
    }),
    EmojisSuggestion.configure({
      HTMLAttributes: {
        class: 'emoji'
      },
      renderHTML({ options, node }) {
        if (node.attrs.emoji) {
          return ['span', options.HTMLAttributes, `${node.attrs.emoji}`];
        } else {
          return [
            'img',
            {
              src: node.attrs.url,
              alt: node.attrs.label,
              class: 'h-[22px] w-[22px] object-contain inline-block'
            }
          ];
        }
      }
    }),
    CustomMention.configure({
      HTMLAttributes: {
        class: 'mention'
      },
      renderHTML({ options, node }) {
        return ['button', options.HTMLAttributes, `${node.attrs.label}`];
      }
    })
  ];
}

export function createEditorConfig({
  element,
  content,
  autofocus,
  placeholder,
  additionalExtensions,
  editorProps,
  onTransaction,
  onEnterPress,
  onBlur,
  onFocus,
  onEscapePress
}: EditorConfigOptions): Partial<EditorOptions> {
  const base = [
    StarterKit.configure({
      gapcursor: false,
      dropcursor: false,
      heading: false,
      orderedList: false,
      bulletList: false,
      blockquote: false,
      link: {
        HTMLAttributes: {
          class: 'rich-link',
          role: 'button',
          rel: 'noopener noreferrer',
        },
        defaultProtocol: 'https'
      }
    }),
    EmojisSuggestion.configure({
      HTMLAttributes: {
        class: 'editor-emoji'
      },
      renderHTML({ options, node }) {
        if (node.attrs.emoji) {
          return ['span', options.HTMLAttributes, `${node.attrs.emoji}`];
        } else {
          return [
            'img',
            {
              src: node.attrs.url,
              alt: node.attrs.label,
              class: 'h-[22px] w-[22px] object-contain inline-block'
            }
          ];
        }
      }
    }),
    CustomMention.configure({
      HTMLAttributes: {
        class: 'editor-mention'
      },
      renderHTML({ options, node }) {
        return [
          'span',
          options.HTMLAttributes,
          `${node.attrs.mentionSuggestionChar}${node.attrs.label}`
        ];
      }
    })
  ];

  const extensions = placeholder ? [placeholder, ...base] : base;

  if (additionalExtensions) extensions.push(...additionalExtensions);

  return {
    element,
    autofocus,
    content,
    extensions: extensions,
    onTransaction,
    onBlur: onBlur ? onBlur : () => { },
    onFocus: onFocus ? onFocus : () => { },
    onUpdate: ({ editor }) => {
      editor.view.dom.scrollTop = editor.view.dom.scrollHeight
    },
    editorProps: editorProps
      ? {
        ...editorProps,
        handleKeyDown: (_, ev) => {
          if (
            ev.key === 'Enter' &&
            !ev.shiftKey &&
            (!editorStore.mentionProps || editorStore.mentionProps.items.length === 0) &&
            (!editorStore.emojiProps || editorStore.emojiProps.items.length === 0)
          ) {
            ev.preventDefault();
            onEnterPress?.();
            return true;
          }

          if (ev.key === 'Escape' && onEscapePress) {
            onEscapePress();
          }


          return false;
        }
      }
      : undefined
  };
}
