import type { APIError } from '$lib/types/errors';
import { errAsync, okAsync, ResultAsync } from 'neverthrow';
import ky, { type Input, type Options } from 'ky';
import type { Category, Channel, Server, Setup, User } from '$lib/types/types';
import type { CreateCategoryType, CreateChannelType, CreateServerType, EditChannelType, PinChannelType } from '$lib/types/schemas';

const client = ky.create({
  prefixUrl: `${import.meta.env.VITE_API_URL}/protected`,
  credentials: 'include',
  retry: 2,
  timeout: 10000
});

export class BackendStore {
  private makeRequest<T>(endpoint: Input, options?: Options): ResultAsync<T, APIError> {
    return ResultAsync.fromPromise(
      client(endpoint, options),
      (error) => ({
        status: 0,
        code: 'NETWORK_ERROR',
        cause: 'Network request failed',
        message: error instanceof Error ? error.message : 'Unknown error',
      })
    ).andThen((res) =>
      ResultAsync.fromPromise(
        res.json(),
        () => ({
          status: res.status,
          code: 'PARSE_ERROR',
          cause: 'Failed to parse response',
          message: 'Invalid JSON response',
        })
      ).andThen((data: unknown) => {
        if (!res.ok) {
          const errorData = data as APIError
          return errAsync({
            status: res.status,
            code: errorData.code || 'API_ERROR',
            cause: errorData.cause || '',
            message: errorData.message || `HTTP ${res.status}`,
          });
        }

        return okAsync(data as T);
      })
    );
  }

  checkIdentity(): ResultAsync<User, APIError> {
    return this.makeRequest<User>('check');
  }

  getSetup(): ResultAsync<Setup, APIError> {
    return this.makeRequest<Setup>('users/setup')
  }

  createServer(body: CreateServerType): ResultAsync<Server, APIError> {
    const formData = new FormData();
    formData.append('name', body.name);
    formData.append('avatar', body.avatar);
    formData.append('crop', JSON.stringify(body.crop));
    formData.append('public', String(body.public));

    if (body.description) formData.append('description', JSON.stringify(body.description));

    return this.makeRequest<Server>('servers', { method: 'post', body: formData })

  }

  createCategory(body: CreateCategoryType): ResultAsync<Category, APIError> {
    return this.makeRequest<Category>('channels/category', { method: 'post', json: body })
  }

  createChannel(body: CreateChannelType): ResultAsync<Channel, APIError> {
    return this.makeRequest<Channel>('channels', { method: 'post', json: body })
  }

  pinChannel(body: PinChannelType): ResultAsync<void, APIError> {
    return this.makeRequest<void>('channels/pin', { method: 'post', json: body })
  }

  deleteChannel(serverID: string, channelID: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`channels/${channelID}`, { method: 'delete', json: { server_id: serverID } })
  }

  deleteCategory(serverID: string, categoryID: string): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`channels/category/${categoryID}`, { method: 'delete', json: { server_id: serverID } })
  }

  editChannel(channelID: string, body: EditChannelType): ResultAsync<void, APIError> {
    return this.makeRequest<void>(`channels/${channelID}`, { method: 'patch', json: body })
  }
}

export const backend = new BackendStore();
