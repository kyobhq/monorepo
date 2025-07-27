// Temporary file while I integrate the design
import People from "ui/icons/People.svelte";
import Globe from "ui/icons/Globe.svelte";

export const SERVERS = [
  {
    id: 1,
    image: 'https://i.pinimg.com/1200x/12/28/d4/1228d4792b5eab4b73a24861d2aed125.jpg',
    href: '/servers/1'
  },
  {
    id: 2,
    image: 'https://i.pinimg.com/736x/63/70/fb/6370fb7ed9fa2946afc67a7002cee8fc.jpg',
    href: '/servers/2'
  },
  {
    id: 3,
    image: 'https://i.pinimg.com/736x/04/ec/41/04ec41b8922ab35fc828dd45368d7a07.jpg',
    href: '/servers/3'
  },
  {
    id: 4,
    image: 'https://i.pinimg.com/736x/9c/f0/10/9cf010b5e285da487ffa74d48c11242e.jpg',
    href: '/servers/4'
  },
  {
    id: 5,
    image: 'https://i.pinimg.com/736x/e8/74/99/e874997bb1c1dc279e6da563463bf94d.jpg',
    href: '/servers/5'
  }
];

export const TABS = [
  { href: '/friends', Icon: People },
  { href: '/servers', Icon: Globe }
];

export const PINNED_CHANNELS = [
  { id: 123, server_id: 1, type: 'textual-e2ee', title: 'The boys', subtitle: 'Private' },
  { id: 182, server_id: 2, type: 'textual', title: 'General', subtitle: 'The Valley' },
  { id: 192, server_id: 3, type: 'textual', title: 'memes', subtitle: 'Psychotics' },
]
