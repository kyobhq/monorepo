import * as v from 'valibot';
import { ChannelTypes } from './types';

export const SignUpSchema = v.object({
  email: v.pipe(
    v.string(),
    v.nonEmpty("Can't be empty."),
    v.email('The email is badly formatted.')
  ),
  username: v.pipe(
    v.string(),
    v.nonEmpty("Can't be empty."),
    v.maxLength(20, 'The length must be equal or below 20 characters.')
  ),
  display_name: v.pipe(
    v.string(),
    v.nonEmpty("Can't be empty."),
    v.maxLength(20, 'The length must be equal or below 20 characters.')
  ),
  password: v.pipe(
    v.string(),
    v.nonEmpty("can't be empty."),
    v.minLength(8, 'Must be at least 8 characters.'),
    v.maxLength(254, 'Must be at most 254 characters.')
  )
});

export const SignInSchema = v.object({
  email: v.pipe(v.string(), v.nonEmpty('Please enter your email.')),
  password: v.pipe(v.string(), v.nonEmpty('Please enter your password.'))
});

export const CreateServerSchema = v.object({
  name: v.pipe(
    v.string(),
    v.maxLength(20, 'Must be at most 20 characters.'),
    v.nonEmpty('Please enter a name for your realm.')
  ),
  description: v.any(),
  avatar: v.pipe(
    v.file('Please select an image file.'),
    v.mimeType(
      ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp', 'image/avif'],
      'Please select a JPEG, PNG, GIF, WEBP or AVIF file.'
    ),
    v.maxSize(1024 * 1024 * 10, 'Please select a file smaller than 10 MB.')
  ),
  public: v.boolean(),
  crop: v.object({
    height: v.number(),
    width: v.number(),
    x: v.number(),
    y: v.number()
  }),
  position: v.number()
});

export interface CreateServerType extends v.InferInput<typeof CreateServerSchema> { }

export const CreateCategorySchema = v.object({
  server_id: v.string(),
  name: v.pipe(
    v.string(),
    v.nonEmpty("Can't be empty."),
    v.maxLength(20, 'Must be at most 20 characters.')
  ),
  position: v.number(),
  users: v.optional(v.array(v.string())),
  roles: v.optional(v.array(v.string())),
  e2ee: v.boolean()
});

export interface CreateCategoryType extends v.InferInput<typeof CreateCategorySchema> { }

export const CreateChannelSchema = v.object({
  position: v.number(),
  category_id: v.string(),
  server_id: v.string(),
  name: v.pipe(
    v.string(),
    v.nonEmpty("Can't be empty."),
    v.maxLength(20, 'Must be at most 20 characters.')
  ),
  description: v.optional(v.string()),
  users: v.optional(v.array(v.string())),
  roles: v.optional(v.array(v.string())),
  type: v.pipe(v.string(), v.enum(ChannelTypes)),
  e2ee: v.boolean()
});

export interface CreateChannelType extends v.InferInput<typeof CreateChannelSchema> { }

export const PinChannelSchema = v.object({
  server_id: v.string(),
  position: v.number()
});

export interface PinChannelType extends v.InferInput<typeof PinChannelSchema> { }

export const EditChannelSchema = v.object({
  server_id: v.string(),
  name: v.pipe(
    v.string(),
    v.nonEmpty("Can't be empty."),
    v.maxLength(20, 'Must be at most 20 characters.')
  ),
  description: v.optional(v.string()),
  users: v.optional(v.array(v.string())),
  roles: v.optional(v.array(v.string()))
});

export interface EditChannelType extends v.InferInput<typeof EditChannelSchema> { }

export const EditCategorySchema = v.object({
  server_id: v.string(),
  name: v.pipe(
    v.string(),
    v.nonEmpty("Can't be empty."),
    v.maxLength(20, 'Must be at most 20 characters.')
  ),
  users: v.optional(v.array(v.string())),
  roles: v.optional(v.array(v.string()))
});

export interface EditCategoryType extends v.InferInput<typeof EditCategorySchema> { }

export const CreateMessageSchema = v.object({
  server_id: v.string(),
  channel_id: v.string(),
  content: v.any(),
  everyone: v.optional(v.boolean()),
  mentions_users: v.optional(v.array(v.string())),
  mentions_roles: v.optional(v.array(v.string())),
  mentions_channels: v.optional(v.array(v.string())),
  attachments: v.optional(v.array(v.pipe(v.file('Please select a valid file.'))))
});

export interface CreateMessageType extends v.InferInput<typeof CreateMessageSchema> { }

export const EditMessageSchema = v.object({
  server_id: v.string(),
  channel_id: v.string(),
  content: v.any(),
  everyone: v.optional(v.boolean()),
  mentions_users: v.optional(v.array(v.string())),
  mentions_roles: v.optional(v.array(v.string())),
  mentions_channels: v.optional(v.array(v.string()))
});

export interface EditMessageType extends v.InferInput<typeof EditMessageSchema> { }

export const DeleteMessageSchema = v.object({
  server_id: v.string(),
  channel_id: v.string(),
  author_id: v.string()
});

export interface DeleteMessageType extends v.InferInput<typeof DeleteMessageSchema> { }

export const EditUserSchema = v.object({
  username: v.string(),
  display_name: v.string(),
  about_me: v.any(),
  links: v.optional(
    v.array(
      v.object({
        id: v.string(),
        label: v.pipe(v.string(), v.maxLength(20, 'Must be at most 20 characters.')),
        url: v.pipe(v.string(), v.url())
      })
    ),
    []
  ),
  facts: v.optional(
    v.array(
      v.object({
        id: v.string(),
        label: v.pipe(v.string()),
        value: v.pipe(v.string(), v.maxLength(20, 'Must be at most 20 characters.'))
      })
    ),
    []
  )
});

export interface EditUserType extends v.InferInput<typeof EditUserSchema> { }

export const EditPasswordSchema = v.pipe(
  v.object({
    current: v.string(),
    new: v.string(),
    confirm: v.string()
  }),
  v.forward(
    v.partialCheck(
      [['new'], ['confirm']],
      (input) => input.new === input.confirm,
      'Passwords do not match.'
    ),
    ['confirm']
  )
);

export interface EditPasswordType extends v.InferInput<typeof EditPasswordSchema> { }

export const EditAvatarSchema = v.object({
  avatar: v.optional(
    v.pipe(
      v.file(),
      v.mimeType(
        ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp', 'image/avif'],
        'Please select a JPEG, PNG, GIF, WEBP or AVIF file.'
      ),
      v.maxSize(1024 * 1024 * 10, 'Please select a file smaller than 10 MB.')
    )
  ),
  banner: v.optional(
    v.pipe(
      v.file(),
      v.mimeType(
        ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp', 'image/avif'],
        'Please select a JPEG, PNG, GIF, WEBP or AVIF file.'
      ),
      v.maxSize(1024 * 1024 * 10, 'Please select a file smaller than 10 MB.')
    )
  )
});

export interface EditAvatarType extends v.InferInput<typeof EditAvatarSchema> { }

export const JoinServerSchema = v.object({
  invite_link: v.pipe(v.string(), v.nonEmpty('Please enter an invite link.'))
});

export interface JoinServerType extends v.InferInput<typeof JoinServerSchema> { }

export const EditServerSchema = v.object({
  name: v.pipe(
    v.string(),
    v.nonEmpty("Can't be empty."),
    v.maxLength(20, 'Must be at most 20 characters.')
  ),
  description: v.any(),
  public: v.boolean()
});

export interface EditServerType extends v.InferInput<typeof EditServerSchema> { }

export const AddEmojisSchema = v.object({
  emojis: v.array(
    v.pipe(
      v.file('Please select an image file.'),
      v.mimeType(
        ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp', 'image/avif'],
        'Please select a JPEG, PNG, GIF, WEBP or AVIF file.'
      ),
      v.maxSize(1024 * 1024 * 1, 'Please select a file smaller than 1 MB.')
    )
  ),
  shortcodes: v.array(v.string())
});

export interface AddEmojisType extends v.InferInput<typeof AddEmojisSchema> { }

export const CreateOrUpdateRoleSchema = v.object({
  id: v.string(),
  name: v.pipe(
    v.string(),
    v.nonEmpty("Can't be empty."),
    v.maxLength(20, 'Must be at most 20 characters.')
  ),
  color: v.string(),
  abilities: v.array(v.string()),
  position: v.number()
});

export interface CreateOrUpdateRoleType extends v.InferInput<typeof CreateOrUpdateRoleSchema> { }

export const AddFriendSchema = v.object({
  friend_username: v.pipe(v.string(), v.nonEmpty("Please enter your friend's username."))
});

export interface AddFriendType extends v.InferInput<typeof AddFriendSchema> { }

export const AcceptFriendSchema = v.object({
  friendship_id: v.string(),
  sender_id: v.string()
});

export interface AcceptFriendType extends v.InferInput<typeof AcceptFriendSchema> { }

export const RemoveFriendSchema = v.object({
  friendship_id: v.string(),
  sender_id: v.string(),
  receiver_id: v.string(),
  channel_id: v.optional(v.string())
});

export interface RemoveFriendType extends v.InferInput<typeof RemoveFriendSchema> { }

export const BanUserSchema = v.object({
  user_id: v.string(),
  reason: v.optional(v.string()),
});

export interface BanUserType extends v.InferInput<typeof BanUserSchema> { }

export const KickUserSchema = v.object({
  user_id: v.string(),
  reason: v.optional(v.string())
});

export interface KickUserType extends v.InferInput<typeof KickUserSchema> { }
