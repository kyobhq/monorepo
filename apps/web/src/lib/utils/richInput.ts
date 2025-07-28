import type { JSONContent } from '@tiptap/core';
import { generateText } from '@tiptap/core';
import StarterKit from '@tiptap/starter-kit';
import { EmojisSuggestion } from 'ui/RichInput/extensions/emojis/emojis';
import { CustomMention } from 'ui/RichInput/extensions/mentions/mentions';

export function trimEmptyNodes(content: JSONContent) {
  if (!content.content || !Array.isArray(content.content)) {
    return content;
  }

  const cleaned = [...content.content];

  while (cleaned.length > 0 && isEmptyParagraph(cleaned[0])) {
    cleaned.shift();
  }

  while (cleaned.length > 0 && isEmptyParagraph(cleaned[cleaned.length - 1])) {
    cleaned.pop();
  }

  return { ...content, content: cleaned };
}

function isEmptyParagraph(node: any) {
  return (
    node.type === 'paragraph' &&
    (!node.content ||
      node.content.length === 0 ||
      (node.content.length === 1 && !node.content[0].text?.trim()))
  );
}

export function extractFirstNParagraphs(htmlString: string, n = 3) {
  const pattern = /<p>.*?<\/p>/gs;
  const matches = htmlString.match(pattern) || [];

  return { paragraphs: matches.slice(0, n).join(''), enoughMatches: matches.length > n };
}

export function generateTextWithExt(content: JSONContent) {
  return generateText(content, [
    StarterKit.configure({
      gapcursor: false,
      dropcursor: false,
      heading: false,
      orderedList: false,
      bulletList: false,
      blockquote: false
    }),
    EmojisSuggestion.configure({
      HTMLAttributes: {
        class: 'emoji'
      },
      renderHTML({ options, node }) {
        return ['span', options.HTMLAttributes, `${node.attrs.emoji}`];
      }
    }),
    CustomMention.configure({
      HTMLAttributes: {
        class: 'mention'
      },
      renderText({ node }) {
        return `<@${node.attrs['user-id']}>`;
      }
    })
  ]);
}
