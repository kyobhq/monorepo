import { PasteRule } from '@tiptap/core';
import Mention from '@tiptap/extension-mention';
import { PluginKey } from '@tiptap/pm/state';
import { serverStore } from 'stores/serverStore.svelte';
import { editorStore } from 'stores/editorStore.svelte';
import { page } from '$app/state';
import { messageStore } from 'stores/messageStore.svelte';
import { backend } from 'stores/backendStore.svelte';
import { logErr } from 'utils/print';

const DEBOUNCE_DELAY = 150;
/** @type {NodeJS.Timeout | null} */
let debounceTimer = null;
/** @type {Map<string, any[]>} */
let searchResultsCache = new Map();

/**
 * @param {string} serverId
 * @param {string} query
 * @returns {Promise<any[]>}
 */
const debouncedSearchMembers = (serverId, query) => {
  return new Promise((resolve) => {
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
    
    const cacheKey = `${serverId}:${query}`;
    if (searchResultsCache.has(cacheKey)) {
      resolve(searchResultsCache.get(cacheKey) || []);
      return;
    }
    
    debounceTimer = setTimeout(async () => {
      try {
        const result = await backend.searchServerMembers(serverId, query);
        result.match(
          (members) => {
            const membersList = members || [];
            searchResultsCache.set(cacheKey, membersList);
            setTimeout(() => searchResultsCache.delete(cacheKey), 30000);
            resolve(membersList);
          },
          (err) => {
            logErr(err);
            resolve([]);
          }
        );
      } catch (error) {
        console.error('Search members error:', error);
        resolve([]);
      }
    }, DEBOUNCE_DELAY);
  });
};

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
      items: async ({ query }) => {
        if (!page.params.server_id) return [];
        const res = [];

        let allMembers = new Set([...Object.values(messageStore.messageAuthors), ...serverStore.getMembers(page.params.server_id)])
        for (const user of allMembers) {
          if (
            user?.username?.toLowerCase().includes(query.toLowerCase()) ||
            user?.display_name?.toLowerCase().includes(query.toLowerCase())
          ) {
            res.push(user);
          }
        }
        
        if (res.length === 0 && query.length >= 2 && query.length <= 20) {
          const apiMembers = await debouncedSearchMembers(page.params.server_id, query);
          if (apiMembers) {
            res.push(...apiMembers);
          }
        }

        return res;
      },
      render: () => {
        return {
          onStart: (props) => {
            editorStore.mentionProps = props;
            editorStore.listOpen = true;
          },
          onUpdate: (props) => {
            editorStore.mentionProps = props;
          },
          onExit: () => {
            editorStore.mentionProps = null;
            editorStore.listOpen = false;
          },
          onKeyDown: (props) => {
            if (props.event.key === 'Escape') {
              editorStore.mentionProps = null;
              editorStore.listOpen = false;
              return true;
            }

            return editorStore.mentionsListEl?.handleKeyDown(props);
          }
        };
      }
    }
  ]
});
