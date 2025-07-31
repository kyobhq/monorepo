import * as v from 'valibot';
import { ChannelTypes } from './types';

export const SignUpSchema = v.object({
  email: v.pipe(
    v.string(),
    v.nonEmpty('Please enter your email.'),
    v.email('The email is badly formatted.')
  ),
  username: v.pipe(
    v.string(),
    v.minLength(1, 'The length must be equal or above 1 character.'),
    v.maxLength(20, 'The length must be equal or below 20 characters.'),
    v.nonEmpty('Please enter a username.')
  ),
  display_name: v.pipe(
    v.string(),
    v.minLength(1, 'The length must be equal or above 1 character.'),
    v.maxLength(20, 'The length must be equal or below 20 characters.'),
    v.nonEmpty('Please enter a display name.')
  ),
  password: v.pipe(
    v.string(),
    v.nonEmpty('Please enter a password.'),
    v.minLength(8, 'This password is too short.'),
    v.maxLength(254, 'This password is too long.')
  )
});

export const SignInSchema = v.object({
  email: v.pipe(v.string(), v.nonEmpty('Please enter your email.')),
  password: v.pipe(v.string(), v.nonEmpty('Please enter your password.'))
});

export const CreateServerSchema = v.object({
  name: v.pipe(
    v.string(),
    v.maxLength(20, 'The length must be equal or below 20 characters.'),
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
  position: v.number(),
})

export interface CreateServerType extends v.InferInput<typeof CreateServerSchema> { }

export const CreateCategorySchema = v.object({
  server_id: v.string(),
  name: v.pipe(
    v.string(),
    v.minLength(1, 'The length must be equal or above 1 character.'),
    v.maxLength(20, 'The length must be equal or below 20 characters.'),
    v.nonEmpty('Please enter a name for your category.')
  ),
  position: v.number(),
  users: v.optional(v.array(v.string())),
  roles: v.optional(v.array(v.string())),
  e2ee: v.boolean()
})

export interface CreateCategoryType extends v.InferInput<typeof CreateCategorySchema> { }

export const CreateChannelSchema = v.object({
  position: v.number(),
  category_id: v.string(),
  server_id: v.string(),
  name: v.pipe(
    v.string(),
    v.minLength(1, 'The length must be equal or above 1 character.'),
    v.maxLength(20, 'The length must be equal or below 20 characters.'),
    v.nonEmpty('Please enter a name for your channel.')
  ),
  description: v.optional(v.string()),
  users: v.optional(v.array(v.string())),
  roles: v.optional(v.array(v.string())),
  type: v.pipe(
    v.string(),
    v.enum(ChannelTypes)
  ),
  e2ee: v.boolean()
})

export interface CreateChannelType extends v.InferInput<typeof CreateChannelSchema> { }

export const PinChannelSchema = v.object({
  server_id: v.string(),
  position: v.number(),
})

export interface PinChannelType extends v.InferInput<typeof PinChannelSchema> { }

export const EditChannelSchema = v.object({
  server_id: v.string(),
  name: v.pipe(
    v.string(),
    v.minLength(1, 'The length must be equal or above 1 character.'),
    v.maxLength(20, 'The length must be equal or below 20 characters.'),
    v.nonEmpty('Please enter a name for your channel.')
  ),
  description: v.optional(v.string()),
  users: v.optional(v.array(v.string())),
  roles: v.optional(v.array(v.string())),
})

export interface EditChannelType extends v.InferInput<typeof EditChannelSchema> { }
