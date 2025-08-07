import type { APIError } from '$lib/types/errors';
import { errAsync, okAsync, ResultAsync } from 'neverthrow';
import ky, { type Input, type Options } from 'ky';
import type {
	Category,
	Channel,
	Message,
	Server,
	ServerInformations,
	Setup,
	User
} from '$lib/types/types';
import type {
	CreateCategoryType,
	CreateChannelType,
	CreateMessageType,
	CreateServerType,
	DeleteMessageType,
	EditAvatarType,
	EditChannelType,
	EditMessageType,
	EditPasswordType,
	EditUserType,
	PinChannelType
} from '$lib/types/schemas';

const client = ky.create({
	prefixUrl: `${import.meta.env.VITE_API_URL}/protected`,
	credentials: 'include',
	retry: 2,
	timeout: 10000,
	throwHttpErrors: false
});

export class BackendStore {
	private makeRequest<T>(endpoint: Input, options?: Options): ResultAsync<T, APIError> {
		return ResultAsync.fromPromise(client(endpoint, options), (error: unknown) => ({
			status: 0,
			code: 'NETWORK_ERROR',
			cause: 'Network request failed',
			message: error instanceof Error ? error.message : 'Unknown error'
		})).andThen((res) =>
			ResultAsync.fromPromise(res.json(), () => ({
				status: res.status,
				code: 'PARSE_ERROR',
				cause: 'Failed to parse response',
				message: 'Invalid JSON response'
			})).andThen((data: unknown) => {
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
			})
		);
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

	deleteChannel(serverID: string, channelID: string): ResultAsync<void, APIError> {
		return this.makeRequest<void>(`channels/${channelID}`, {
			method: 'delete',
			json: { server_id: serverID }
		});
	}

	deleteCategory(serverID: string, categoryID: string): ResultAsync<void, APIError> {
		return this.makeRequest<void>(`channels/category/${categoryID}`, {
			method: 'delete',
			json: { server_id: serverID }
		});
	}

	editChannel(channelID: string, body: EditChannelType): ResultAsync<void, APIError> {
		return this.makeRequest<void>(`channels/${channelID}`, { method: 'patch', json: body });
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

	getMessages(channelID: string): ResultAsync<Message[], APIError> {
		return this.makeRequest<Message[]>(`messages/${channelID}`);
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

	getUserProfile(userID: string): ResultAsync<User, APIError> {
		return this.makeRequest<User>(`users/${userID}`);
	}

	getInviteLink(serverID: string): ResultAsync<string, APIError> {
		return this.makeRequest<string>(`servers/${serverID}/invite`, { method: 'post' });
	}

	joinServerWithInvite(inviteID: string): ResultAsync<Server, APIError> {
		return this.makeRequest<Server>(`servers/join`, {
			method: 'post',
			json: { invite_id: inviteID }
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
}

export const backend = new BackendStore();
