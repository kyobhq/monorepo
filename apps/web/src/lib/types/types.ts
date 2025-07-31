export const ChannelTypes = {
  Textual: 'textual',
  TextualE2EE: 'textual-e2ee',
  Voice: 'voice',
  Dm: 'dm'
} as const;
export type ChannelTypes = (typeof ChannelTypes)[keyof typeof ChannelTypes];

export const ABILITIES = ['ADMIN', 'MANAGE_CHANNELS', 'MANAGE_ROLES', 'MANAGE_SERVER', 'MANAGE_EXPRESSIONS', 'CHANGE_NICKNAME', 'MANAGE_NICKNAMES', 'BAN', 'KICK', 'MUTE', 'ATTACH_FILES', 'MANAGE_MESSAGES'] as const
export type AbilitiesType = typeof ABILITIES[number]

export interface Channel {
  id: string;
  position: number;
  server_id: string;
  category_id: string;
  name: string;
  description: string;
  type: ChannelTypes;
  unread: boolean;
  last_message_sent?: string;
  last_message_read?: string;
  last_mentions?: string[];
  messages?: Message[];
  users?: string[];
  roles?: string[];
  voice_users: {
    user_id: string;
    deafen: boolean;
    mute: boolean;
  }[];
}

export interface Server {
  id: string;
  position: number;
  owner_id: string;
  name: string;
  avatar: string;
  banner: string;
  description?: any;
  main_color?: string;
  categories: Record<string, Category>;
  public: boolean;
}

export interface Role {
  id: string;
  idx: number;
  name: string
  color: string
  abilities: string[]
  members: string[]
}

export interface Fact {
  id: string;
  label: string;
  value: string;
}

export interface Link {
  id: string;
  label: string;
  url: string;
}

export interface User {
  id: string;
  email: string;
  username: string;
  display_name: string;
  avatar: string;
  banner: string;
  main_color?: string;
  about?: any;
  facts: Fact[];
  links: Link[];
}

export interface Member extends Partial<User> {
  roles: string[];
}

export interface Friend extends Partial<User> {
  channel_id?: string;
  friendship_id: string;
  accepted: boolean;
  sender: boolean;
}

export interface Setup {
  servers: Record<string, Server>;
  emojis: Emoji[];
  friends: Friend[];
}

export interface DefaultResponse {
  message: string;
}

export interface Attachment {
  id: string;
  url: string;
  file_name: string;
  file_size: string;
  type: string;
}

export interface Message {
  id: string;
  author_id: string;
  server_id: string;
  channel_id: string;
  content: any;
  everyone: boolean;
  mentions_users: string[];
  mentions_channels: string[];
  attachments: Attachment[];
  updated_at: string;
  created_at: string;
}

export interface LastState {
  channel_ids: string[];
  last_message_ids: string[];
  mentions_ids: string[][];
}

export interface Emoji {
  id: string;
  url: string;
  shortcode: string;
}

export interface ContextMenuTarget {
  name: string
  author?: string
}


export interface Category {
  id: string;
  position: number;
  server_id: string;
  name: string;
  users?: string[];
  roles?: string[];
  e2ee: boolean;
  channels: Record<string, Channel>;
}
