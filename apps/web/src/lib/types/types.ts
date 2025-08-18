import type { Abilities } from "$lib/constants/permissions";

export const ChannelTypes = {
  Textual: 'textual',
  TextualE2EE: 'textual-e2ee',
  Voice: 'voice',
  Dm: 'dm'
} as const;
export type ChannelTypes = (typeof ChannelTypes)[keyof typeof ChannelTypes];

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
  voice_users?: {
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
  members: Member[];
  user_roles: string[];
  roles: Role[];
  invites: Invite[];
}

export interface Invite {
  id: string;
  creator: Partial<User>;
  server_id: string;
  invite_id: string;
  expire_at: string;
}

export interface Role {
  id: string;
  position: number;
  name: string;
  color: string;
  abilities: Abilities[];
  members: string[];
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
  about_me?: any;
  facts: Fact[];
  links: Link[];
  updated_at: string;
  created_at: string;
}

export interface Member extends Partial<User> {
  status: string;
  joined_kyob: string;
  joined_server: string;
  roles: string[];
}

export interface ServerInformations {
  member_count: number;
  members: Member[];
  roles: Role[];
  user_roles: string[];
  invites: Invite[];
}

export interface Friend extends Partial<User> {
  channel_id?: string;
  friendship_id: string;
  friendship_sender_id: string;
  accepted: boolean;
  status: string;
}

export interface Setup {
  user: User;
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
  author: {
    id: string;
    avatar: string;
    display_name: string;
  };
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
  name: string;
  author?: string;
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
