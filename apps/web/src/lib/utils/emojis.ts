import { userStore } from 'stores/userStore.svelte';
import emojis from 'emojibase-data/en/compact.json';
import type { Emoji } from '$lib/types/types';

interface EmojiMatch {
  emoji: Emoji | any;
  score: number;
  type: 'personal' | 'default';
}

export function searchEmoji(query: string, limit: number = 8) {
  const lowerQuery = query.toLocaleLowerCase().trim();

  if (!lowerQuery) return [];

  const matches: EmojiMatch[] = [];

  if (userStore.emojis) {
    for (const emoji of userStore.emojis) {
      const score = scorePersonalEmoji(emoji, lowerQuery);
      if (score > 0) {
        matches.push({ emoji, score: score + 100, type: 'personal' });
      }
    }
  }

  for (const emoji of emojis) {
    const score = scoreDefaultEmoji(emoji, lowerQuery);
    if (score > 0) {
      matches.push({ emoji, score, type: 'default' });
    }
  }

  return matches.sort((a, b) => b.score - a.score).slice(0, limit);
}

function scorePersonalEmoji(emoji: Emoji, query: string): number {
  const shortcode = emoji.shortcode;

  if (shortcode === query) return 100;
  if (shortcode.startsWith(query)) return 80;
  if (shortcode.includes(query)) return 40;

  return 0;
}

function scoreDefaultEmoji(emoji: { label: string; tags?: string[] }, query: string): number {
  const label = emoji.label.toLowerCase().replaceAll(' ', '_');
  const tags = emoji.tags;

  if (!label) return 0;

  if (label === query) return 100;
  if (label.startsWith(query)) return 80;
  if (label.includes(query)) return 70;
  if (tags?.includes(query)) return 60;

  const tagMatch = tags?.some((tag) => tag.includes(query));
  if (tagMatch) return 30;

  return 0;
}

export function transformShortcode(input: string) {
  return input
    .toLowerCase()
    .replace(/\s+/g, '_')
    .replace(/[^a-z0-9_]/g, '')
    .slice(0, 20);
}
