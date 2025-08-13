export const ABILITIES = {
  VIEW_CHANNELS: 'VIEW_CHANNELS',
  MANAGE_CHANNELS: 'MANAGE_CHANNELS',
  MANAGE_ROLES: 'MANAGE_ROLES',
  MANAGE_SERVER: 'MANAGE_SERVER',
  CREATE_INVITE: 'CREATE_INVITE',
  CHANGE_NICKNAME: 'CHANGE_NICKNAME',
  MANAGE_NICKNAMES: 'MANAGE_NICKNAMES',
  TIMEOUT_MEMBERS: 'TIMEOUT_MEMBERS',
  KICK_MEMBERS: 'KICK_MEMBERS',
  BAN_MEMBERS: 'BAN_MEMBERS',
  SEND_MESSAGES: 'SEND_MESSAGES',
  ATTACH_FILES: 'ATTACH_FILES',
  ADD_REACTIONS: 'ADD_REACTIONS',
  USE_PERSONAL_EMOJIS: 'USE_PERSONAL_EMOJIS',
  MENTION_EVERYONE: 'MENTION_EVERYONE',
  MANAGE_MESSAGES: 'MANAGE_MESSAGES',
  CONNECT: 'CONNECT',
  SPEAK: 'SPEAK',
  VIDEO: 'VIDEO',
  MUTE_MEMBERS: 'MUTE_MEMBERS',
  DEAFEN_MEMBERS: 'DEAFEN_MEMBERS',
  MOVE_MEMBERS: 'MOVE_MEMBERS',
  ADMINISTRATOR: 'ADMINISTRATOR',
  OWNER: 'OWNER'
} as const;

export type Abilities = (typeof ABILITIES)[keyof typeof ABILITIES];

export const ABILITIES_CONTENT = [
  {
    code: ABILITIES.VIEW_CHANNELS,
    label: 'View Channels',
    description: 'Allows members to view channels by default (excluding private channels).'
  },
  {
    code: ABILITIES.MANAGE_CHANNELS,
    label: 'Manage Channels',
    description: 'Allows members to create, edit, or delete channels.'
  },
  {
    code: ABILITIES.MANAGE_ROLES,
    label: 'Manage Roles',
    description:
      'Allow members to create new roles and edit or delete roles lower than their highest role.'
  },
  {
    code: ABILITIES.MANAGE_SERVER,
    label: 'Manage Server',
    description:
      "Allows members to change this server's name, description, avatar, banner, view all invites, and members."
  },
  {
    code: ABILITIES.CREATE_INVITE,
    label: 'Create Invite',
    description: 'Allows members to invite new people to this server.'
  },
  {
    code: ABILITIES.CHANGE_NICKNAME,
    label: 'Change Nickname',
    description: 'Allows members to change their own nickname.'
  },
  {
    code: ABILITIES.MANAGE_NICKNAMES,
    label: 'Manage Nicknames',
    description: 'Allows members to change the nicknames of other members.'
  },
  {
    code: ABILITIES.TIMEOUT_MEMBERS,
    label: 'Timeout Members',
    description:
      "Allows members to timeout other members. When timed out the user won't be able to send any message, react or speak in voice channels."
  },
  {
    code: ABILITIES.KICK_MEMBERS,
    label: 'Kick Members',
    description:
      'Allows members to remove other members from this server. Kicked members can join again if they have access to another invite.'
  },
  {
    code: ABILITIES.BAN_MEMBERS,
    label: 'Ban Members',
    description:
      'Allows members to permanently ban and delete the message history of other members on this server.'
  },
  {
    code: ABILITIES.SEND_MESSAGES,
    label: 'Send Messages',
    description: 'Allows members to send messages in text channels.'
  },
  {
    code: ABILITIES.ATTACH_FILES,
    label: 'Attach Files',
    description: 'Allows members to upload files or media in text channels.'
  },
  {
    code: ABILITIES.ADD_REACTIONS,
    label: 'Add Reactions',
    description:
      'Allows members to add new emoji reactions to a message or use any already existing reaction.'
  },
  {
    code: ABILITIES.USE_PERSONAL_EMOJIS,
    label: 'Use Personal Emojis',
    description: 'Allows members to use their own emojis.'
  },
  {
    code: ABILITIES.MENTION_EVERYONE,
    label: 'Mention @everyone and all Roles',
    description:
      'Allows members to use @everyone (everyone on the server). They can also @mention all roles, even if the role\'s "Allow anyone to mention this role" permission is disabled.'
  },
  {
    code: ABILITIES.MANAGE_MESSAGES,
    label: 'Manage Messages',
    description: 'Allows members to delete messages by other members.'
  },
  {
    code: ABILITIES.CONNECT,
    label: 'Connect',
    description: 'Allows members to join voice channels.'
  },
  {
    code: ABILITIES.SPEAK,
    label: 'Speak',
    description:
      'Allows members to speak in voice channels. If this permission is disabled they\'ll be muted by default until somebody with the "Mute Members" permission un-mute them.'
  },
  {
    code: ABILITIES.VIDEO,
    label: 'Video',
    description:
      'Allows members to share their video, screen share, or stream a game in this server.'
  },
  {
    code: ABILITIES.MUTE_MEMBERS,
    label: 'Mute Members',
    description: 'Allows members to mute other members in voice channels for everyone.'
  },
  {
    code: ABILITIES.DEAFEN_MEMBERS,
    label: 'Deafen Members',
    description:
      "Allows members to deafen other members in voice channels, which means they won't be able to hear others."
  },
  {
    code: ABILITIES.MOVE_MEMBERS,
    label: 'Move Members',
    description:
      'Allows members to disconnect or move other members between voice channels that the member with this permission has access to.'
  },
  {
    code: ABILITIES.ADMINISTRATOR,
    label: 'Administrator',
    description:
      "Members with this permission have every rights on this server. They'll be able to bypass any restrictions put on channels. This is a dangerous permission to grant be careful."
  }
] as const;

export type AbilitiesContent = (typeof ABILITIES_CONTENT)[number];
