import * as v from 'valibot';

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
})

export interface CreateServerType extends v.InferInput<typeof CreateServerSchema> { }
