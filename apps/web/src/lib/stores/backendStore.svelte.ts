import type { APIError } from '$lib/types/errors';
import { errAsync, okAsync, ResultAsync } from 'neverthrow';
import ky, { type Input, type Options } from 'ky';
import type {
  Category,
  Channel,
  Emoji,
  Friend,
  Member,
  Message,
  Role,
  Server,
  ServerInformations,
  Setup,
  User
} from '$lib/types/types';
import type {
  AcceptFriendType,
  AddEmojisType,
  AddFriendType,
  BanUserType,
  CreateCategoryType,
  CreateChannelType,
  CreateMessageType,
  CreateOrUpdateRoleType,
  CreateServerType,
  DeleteMessageType,
  EditAvatarType,
  EditCategoryType,
  EditChannelType,
  EditMessageType,
  EditPasswordType,
  EditServerType,
  EditUserType,
  KickUserType,
  PinChannelType,
  RemoveFriendType
} from '$lib/types/schemas';
import { channelStore } from './channelStore.svelte';
import { SvelteMap } from 'svelte/reactivity';

const client = ky.create({
  prefixUrl: `${import.meta.env.VITE_API_URL}/protected`,
  credentials: 'include',
  retry: 2,
  timeout: 10000,
  throwHttpErrors: false
});

export class BackendStore {
  private pendingRequests = new SvelteMap<string, AbortController>();

  private makeRequest<T>(endpoint: Input, options?: Options): ResultAsync<T, APIError> {
    const method = options?.method || 'get';
    const bodyKey = options?.json
      ? JSON.stringify(options.json)
      : options?.body
        ? String(options.body)
        : '';
    const requestKey = `${method.toUpperCase()}:${String(endpoint)}:${bodyKey}`;

    if (this.pendingRequests.has(requestKey)) {
      this.pendingRequests.get(requestKey)?.abort();
    }

    const controller = new AbortController();
    this.pendingRequests.set(requestKey, controller);

    const shouldTrace =
      import.meta.env.DEV &&
      typeof localStorage !== 'undefined' &&
      localStorage.getItem('traceApi') === '1';
    const startAt = shouldTrace ? performance.now() : 0;

    return ResultAsync.fromPromise(
      client(endpoint, { ...options, signal: controller.signal }),
      (error: unknown) => {
        this.pendingRequests.delete(requestKey);

        if (error instanceof Error && error.name === 'AbortError') {
          return {
            status: 0,
            code: 'REQUEST_ABORTED',
            cause: 'Request was aborted',
            message: 'Request cancelled by newer request'
          };
        }

        return {
          status: 0,
          code: 'NETWORK_ERROR',
          cause: 'Network request failed',
          message: error instanceof Error ? error.message : 'Unknown error'
        };
      }
    ).andThen((res) => {
      this.pendingRequests.delete(requestKey);
      const afterFetchAt = shouldTrace ? performance.now() : 0;
      return ResultAsync.fromPromise(res.json(), () => ({
        status: res.status,
        code: 'PARSE_ERROR',
        cause: 'Failed to parse response',
        message: 'Invalid JSON response'
      })).andThen((data: unknown) => {
        if (shouldTrace) {
          const afterJsonAt = performance.now();

          console.log('[API]', String(endpoint), {
            networkMs: +(afterFetchAt - startAt).toFixed(2),
            parseMs: +(afterJsonAt - afterFetchAt).toFixed(2),
            totalMs: +(afterJsonAt - startAt).toFixed(2)
          });
        }
        if (!res.ok) {
          const errorData = data as APIError;
          return errAsync({
            status: res.status,
            code: errorData.code || 'API_ERROR',
            cause: errorData.cause || '',
            message: errorData.message || `HTTP ${res.status}`
          });
        }

        return okAsync(data as T);
      });
    });
  }

  getSetup(): ResultAsync<Setup, APIError> {
    return this.makeRequest<Setup>('users/setup');
  }

  createServer(body: CreateServerType): ResultAsync<Server, APIError> {
    const formData = new FormData();
    formData.append('name', body.name);
    formData.append('avatar', body.avatar);
    formData.append('crop', JSON.stringify(body.crop));
    formData.append('public', String(body.public));

    if (body.description) formData.append('description', JSON.stringify(body.description));

    return this.makeRequest<Server>('servers', { method: 'post', body: formData });
  }

  createCategory(body: CreateCategoryType): ResultAsync<Category, APIError> {
    return this.makeRequest<Category>('channels/category', { method: 'post', json: body });
  }

  createChannel(body: CreateChannelType): ResultAsync<Channel, APIError> {
    return this.makeRequest<Channel>('channels', { method: 'post', json: body });
  }

  pinChannel(body: PinChannelType): ResultAsync<void, APIError> {
    return this.makeRequest<void>('channels/pin', { method: 'post', json: body });
  }

  deleteChannel(
    serverID: string,
    categoryID: string,
    channelID: string
  ): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`channels/${channelID}`, {
      method: 'delete',
      json: { server_id: serverID, category_id: categoryID }
    });
  }

  deleteServer(serverID: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`servers/${serverID}`, { method: 'delete' });
  }

  deleteCategory(serverID: string, categoryID: string): ResultAsync<void, APIError> {
    const channels = channelStore.getCategoryChannels(serverID, categoryID);
    return this.makeRequest<void>(`channels/category/${categoryID}`, {
      method: 'delete',
      json: { server_id: serverID, channels_ids: channels.map((chan) => chan.id) }
    });
  }

  editChannel(channelID: string, body: EditChannelType): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`channels/${channelID}`, { method: 'patch', json: body });
  }

  editCategory(categoryID: string, body: EditCategoryType): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`channels/category/${categoryID}`, {
      method: 'patch',
      json: body
    });
  }

  createMessage(body: CreateMessageType): ResultAsync<void, APIError> {
    const formData = new FormData();
    formData.append('server_id', body.server_id);
    formData.append('channel_id', body.channel_id);
    formData.append('content', JSON.stringify(body.content));
    formData.append('everyone', body.everyone ? 'true' : 'false');
    body.mentions_users?.forEach((user) => formData.append('mentions_users[]', user));
    body.mentions_channels?.forEach((channel) => formData.append('mentions_channels[]', channel));
    body.mentions_roles?.forEach((role) => formData.append('mentions_roles[]', role));
    body.attachments?.forEach((attachment) => formData.append('attachments[]', attachment));

    return this.makeRequest(`messages`, { method: 'post', body: formData });
  }

  getServerInformations(serverID: string): ResultAsync<ServerInformations, APIError> {
    return this.makeRequest<ServerInformations>(`servers/${serverID}`);
  }

  getMessages(serverID: string, channelID: string): ResultAsync<Message[], APIError> {
    return this.makeRequest<Message[]>(`messages/${serverID}/${channelID}`);
  }

  editMessage(messageID: string, body: EditMessageType): ResultAsync<void, APIError> {
    return this.makeRequest(`messages/${messageID}`, { method: 'patch', json: body });
  }

  deleteMessage(messageID: string, body: DeleteMessageType): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`messages/${messageID}`, { method: 'delete', json: body });
  }

  updateProfile(body: EditUserType): ResultAsync<void, APIError> {
    return this.makeRequest<void>('users/profile', { method: 'patch', json: body });
  }

  updatePassword(body: EditPasswordType): ResultAsync<void, APIError> {
    return this.makeRequest<void>('users/password', { method: 'patch', json: body });
  }

  updateAvatarAndBanner(
    cropAvatarPixels: any,
    cropBannerPixels: any,
    body: EditAvatarType
  ): ResultAsync<{ avatar: string; banner: string }, APIError> {
    const formData = new FormData();
    if (body.avatar) formData.append('avatar', body.avatar);
    if (body.banner) formData.append('banner', body.banner);
    if (cropAvatarPixels) formData.append('crop_avatar', JSON.stringify(cropAvatarPixels));
    if (cropBannerPixels) formData.append('crop_banner', JSON.stringify(cropBannerPixels));

    return this.makeRequest<{ avatar: string; banner: string }>('users/avatar', {
      method: 'patch',
      body: formData
    });
  }

  updateServerAvatarAndBanner(
    serverID: string,
    cropAvatarPixels: any,
    cropBannerPixels: any,
    body: EditAvatarType
  ): ResultAsync<{ avatar: string; banner: string }, APIError> {
    const formData = new FormData();
    if (body.avatar) formData.append('avatar', body.avatar);
    if (body.banner) formData.append('banner', body.banner);
    if (cropAvatarPixels) formData.append('crop_avatar', JSON.stringify(cropAvatarPixels));
    if (cropBannerPixels) formData.append('crop_banner', JSON.stringify(cropBannerPixels));

    return this.makeRequest<{ avatar: string; banner: string }>(`servers/${serverID}/avatar`, {
      method: 'patch',
      body: formData
    });
  }

  getUserProfile(userID: string): ResultAsync<User, APIError> {
    return this.makeRequest<User>(`users/${userID}`);
  }

  getInviteLink(serverID: string): ResultAsync<string, APIError> {
    return this.makeRequest<string>(`servers/${serverID}/invite`, { method: 'post' });
  }

  joinServerWithInvite(inviteID: string, position: number): ResultAsync<Server, APIError> {
    return this.makeRequest<Server>(`servers/join`, {
      method: 'post',
      json: { invite_id: inviteID, position: position }
    });
  }

  joinPublicServer(serverID: string): ResultAsync<Server, APIError> {
    return this.makeRequest<Server>(`servers/join`, {
      method: 'post',
      json: { server_id: serverID }
    });
  }

  leaveServer(serverID: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`servers/${serverID}/leave`, { method: 'post' });
  }

  editServerProfile(serverID: string, body: EditServerType): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`servers/${serverID}/profile`, { method: 'patch', json: body });
  }

  uploadEmojis(body: AddEmojisType): ResultAsync<Emoji[], APIError> {
    const formData = new FormData();
    for (let i = 0; i < body.emojis.length; ++i) {
      formData.append('emojis[]', body.emojis[i]);
      formData.append('shortcodes[]', body.shortcodes[i]);
    }

    return this.makeRequest<Emoji[]>(`users/emojis`, { method: 'post', body: formData });
  }

  updateEmoji(id: string, shortcode: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`users/emojis/${id}`, { method: 'patch', json: { shortcode } });
  }

  deleteEmoji(id: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`users/emojis/${id}`, { method: 'delete' });
  }

  createOrUpdateRole(serverID: string, body: CreateOrUpdateRoleType): ResultAsync<Role, APIError> {
    return this.makeRequest<Role>(`roles`, {
      method: 'post',
      json: { ...body, server_id: serverID }
    });
  }

  deleteRole(serverID: string, roleID: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`roles`, {
      method: 'delete',
      json: { server_id: serverID, role_id: roleID }
    });
  }

  addRoleMember(serverID: string, roleID: string, userID: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`roles/add_member`, {
      method: 'patch',
      json: { server_id: serverID, role_id: roleID, user_id: userID }
    });
  }

  removeRoleMember(serverID: string, roleID: string, userID: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`roles/remove_member`, {
      method: 'patch',
      json: { server_id: serverID, role_id: roleID, user_id: userID }
    });
  }

  moveRole(
    serverID: string,
    targetRoleID: string,
    movedRoleID: string,
    from: number,
    to: number
  ): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`roles/move`, {
      method: 'patch',
      json: {
        server_id: serverID,
        target_role_id: targetRoleID,
        moved_role_id: movedRoleID,
        from: from,
        to: to
      }
    });
  }

  sendFriendRequest(body: AddFriendType): ResultAsync<Friend, APIError> {
    return this.makeRequest<Friend>('friends', { method: 'post', json: body });
  }

  acceptFriendRequest(body: AcceptFriendType): ResultAsync<void, APIError> {
    return this.makeRequest<void>('friends', { method: 'patch', json: body });
  }

  removeFriend(body: RemoveFriendType): ResultAsync<void, APIError> {
    return this.makeRequest<void>('friends', { method: 'delete', json: body });
  }

  deleteAccount(): ResultAsync<void, APIError> {
    return this.makeRequest<void>('users', { method: 'delete' });
  }

  banUser(serverID: string, body: BanUserType): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`servers/${serverID}/ban`, { method: 'post', json: body });
  }

  unbanUser(serverID: string, userID: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`servers/${serverID}/unban/${userID}`, { method: 'post' });
  }

  kickUser(serverID: string, body: KickUserType): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`servers/${serverID}/kick`, { method: 'post', json: body });
  }

  getBannedMembers(serverID: string): ResultAsync<Member[], APIError> {
    return this.makeRequest<Member[]>(`servers/${serverID}/bans`);
  }
}

export const backend = new BackendStore();
