import { PasteRule } from '@tiptap/core';
import Mention from '@tiptap/extension-mention';
import { PluginKey } from '@tiptap/pm/state';
import { serverStore } from 'stores/serverStore.svelte';
import { editorStore } from 'stores/editorStore.svelte';
import { page } from '$app/state';

const MentionExtended = Mention.extend({
  addPasteRules() {
    return [
      new PasteRule({
        find: /<@(\d+)>/g,
        handler: ({ state, range, match }) => {
          if (!page.params.server_id) return;

          const userId = match[1];
          const user = serverStore.getMember(page.params.server_id, userId);

          const attributes = {
            'user-id': userId,
            label: user?.display_name || 'Unknown User',
            mentionSuggestionChar: '@'
          };

          const { tr } = state;
          tr.replaceWith(range.from, range.to, this.type.create(attributes));
        }
      })
    ];
  },

  addAttributes() {
    return {
      'user-id': {
        default: null
      },
      label: {
        default: null
      },
      mentionSuggestionChar: {
        default: '@'
      }
    };
  },

  addStorage() {
    return {
      mentionProps: null,
      mentionsListEl: null
    };
  }
});

export const CustomMention = MentionExtended.configure({
  renderText({ node }) {
    return `<@${node.attrs['user-id']}>`;
  },

  suggestions: [
    {
      char: '@',
      pluginKey: new PluginKey('at'),
      items: ({ query }) => {
        if (!page.params.server_id) return [];
        const res = [];

        for (const user of serverStore.getMembers(page.params.server_id)) {
          if (
            user?.username?.toLowerCase().includes(query.toLowerCase()) ||
            user?.display_name?.toLowerCase().includes(query.toLowerCase())
          ) {
            res.push(user);
          }
        }

        return res;
      },
      render: () => {
        return {
          onStart: (props) => {
            editorStore.mentionProps = props;
          },
          onUpdate: (props) => {
            editorStore.mentionProps = props;
          },
          onExit: () => {
            editorStore.mentionProps = null;
          },
          onKeyDown: (props) => {
            if (props.event.key === 'Escape') {
              editorStore.mentionProps = null;
              return true;
            }

            return editorStore.mentionsListEl?.handleKeyDown(props);
          }
        };
      }
    }
  ]
});
