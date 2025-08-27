import { WSMessageSchema, type WSMessage } from '$lib/gen/types_pb';
import type {
	Category,
	Channel,
	ChannelTypes,
	Friend,
	Member,
	Message,
	Role
} from '$lib/types/types';
import { fromBinary } from '@bufbuild/protobuf';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { channelStore } from './channelStore.svelte';
import { page } from '$app/state';
import { serverStore } from './serverStore.svelte';
import { userStore } from './userStore.svelte';
import { categoryStore } from './categoryStore.svelte';
import { goto } from '$app/navigation';
import { coreStore } from './coreStore.svelte';
import type {
	ChangeStatus,
	DeleteChatMessage,
	EditChatMessage,
	NewChatMessage,
	StartCategory,
	StartChannel,
	KillCategory,
	KillChannel,
	CreateOrEditRole,
	AddRoleMember,
	RemoveRoleMember,
	RemoveRole,
	MoveRole,
	FriendRequest,
	AcceptFriendRequest,
	RemoveFriend,
	AccountDeletion,
	AvatarServerChange,
	ProfileServerChange,
	EditChannel,
	EditCategory,
	KillServer,
	LeaveServer,
	BanUser,
	KickUser
} from '$lib/gen/types_pb';

export class WebsocketStore {
	wsConn = $state<WebSocket>();

	init(userID: string) {
		const ws = new WebSocket(`ws://localhost:8080/api/protected/ws/${userID}`);
		if (!ws) return;

		this.wsConn = ws;

		ws.onopen = () => {
			window.setInterval(() => {
				ws.send('heartbeat');
			}, 10 * 1000);
		};

		ws.onmessage = async (event) => {
			if (event.data === 'heartbeat') return;

			const arrayBuffer = await event.data.arrayBuffer();
			const uint8Array = new Uint8Array(arrayBuffer);
			const wsMess = fromBinary(WSMessageSchema, uint8Array, {
				readUnknownFields: false
			});

			await this.handleMessage(wsMess);
		};
	}

	private async handleMessage(wsMess: WSMessage) {
		const { content } = wsMess;
		if (!content.case || !content.value) return;

		const handlers: Record<string, () => void> = {
			userChangeStatus: () => this.handleUserStatusChange(content.value as ChangeStatus),
			newChatMessage: () => this.handleNewMessage(content.value as NewChatMessage),
			deleteChatMessage: () => this.handleDeleteMessage(content.value as DeleteChatMessage),
			editChatMessage: () => this.handleEditMessage(content.value as EditChatMessage),
			startCategory: () => this.handleCategoryStart(content.value as StartCategory),
			startChannel: () => this.handleChannelStart(content.value as StartChannel),
			killCategory: () => this.handleCategoryDelete(content.value as KillCategory),
			killChannel: () => this.handleChannelDelete(content.value as KillChannel),
			createOrEditRole: () => this.handleRoleCreateOrEdit(content.value as CreateOrEditRole),
			addRoleMember: () => this.handleAddRoleMember(content.value as AddRoleMember),
			removeRoleMember: () => this.handleRemoveRoleMember(content.value as RemoveRoleMember),
			removeRole: () => this.handleRoleDelete(content.value as RemoveRole),
			moveRole: () => this.handleRoleMove(content.value as MoveRole),
			friendRequest: () => this.handleFriendRequest(content.value as FriendRequest),
			acceptFriendRequest: () => this.handleAcceptFriend(content.value as AcceptFriendRequest),
			removeFriend: () => this.handleRemoveFriend(content.value as RemoveFriend),
			accountDeletion: () => this.handleAccountDeletion(content.value as AccountDeletion),
			avatarServerChange: () => this.handleServerAvatarChange(content.value as AvatarServerChange),
			profileServerChange: () =>
				this.handleServerProfileChange(content.value as ProfileServerChange),
			editChannel: () => this.handleChannelEdit(content.value as EditChannel),
			editCategory: () => this.handleCategoryEdit(content.value as EditCategory),
			killServer: () => this.handleServerDelete(content.value as KillServer),
			leaveServer: () => this.handleServerLeave(content.value as LeaveServer),
			banUser: () => this.handleUserBan(content.value as BanUser),
			kickUser: () => this.handleUserKick(content.value as KickUser)
		};

		const handler = handlers[content.case];
		if (handler) handler();
	}

	private handleUserStatusChange(value: ChangeStatus) {
		if (!value?.user || value.user.id === userStore.user?.id) return;

		if (value.status === 'offline') {
			serverStore.setMemberOffline(value.serverId, value.user.id);
			userStore.setFriendStatus(value.user.id, value.status);
			return;
		}

		if (value.type === 'connect') {
			serverStore.setMemberOnline(value.serverId, value.user.id, value.status);
			userStore.setFriendStatus(value.user.id, value.status);
		}

		if (value.type === 'join') {
			const member: Member = {
				id: value.user.id,
				display_name: value.user.displayName,
				avatar: value.user.avatar,
				status: value.status,
				roles: [],
				joined_server: Date.now().toString(),
				joined_kyob: userStore.user!.created_at
			};
			serverStore.addMember(value.serverId, member);
		}
	}

	private handleNewMessage(value: NewChatMessage) {
		if (!value.message) return;

		const msg = value.message;
		const newMessage: Message = {
			id: msg.id,
			author: {
				id: msg.author!.id,
				avatar: msg.author!.avatar,
				display_name: msg.author!.displayName,
				status: 'online',
				joined_kyob: '',
				joined_server: '',
				roles: []
			},
			server_id: msg.serverId,
			channel_id: msg.channelId,
			content: JSON.parse(new TextDecoder().decode(msg.content)),
			everyone: msg.everyone,
			mentions_users: msg.mentionsUsers,
			mentions_channels: msg.mentionsChannels,
			attachments: this.parseAttachments(msg.attachments),
			updated_at: timestampDate(msg.createdAt!).toISOString(),
			created_at: timestampDate(msg.createdAt!).toISOString()
		};

		if (msg.serverId === 'global') {
			channelStore.addMessageDM(msg.channelId, newMessage);
		} else {
			channelStore.addMessage(msg.serverId, msg.channelId, newMessage);
		}
	}

	private handleDeleteMessage(value: DeleteChatMessage) {
		if (!value.message) return;
		channelStore.deleteMessage(value.message.channelId, value.message.id);
	}

	private handleEditMessage(value: EditChatMessage) {
		if (!value.message) return;

		const msg = value.message;
		const editMessage: Partial<Message> = {
			id: msg.id,
			everyone: msg.everyone,
			mentions_users: msg.mentionsUsers,
			mentions_channels: msg.mentionsChannels,
			content: JSON.parse(new TextDecoder().decode(msg.content)),
			updated_at: timestampDate(msg.updatedAt!).toISOString()
		};

		channelStore.editMessage(msg.channelId, editMessage);
	}

	private handleCategoryStart(value: StartCategory) {
		if (!value.category) return;

		const cat = value.category;
		const newCategory: Category = {
			id: cat.id,
			server_id: cat.serverId,
			name: cat.name,
			position: cat.position,
			users: cat.users,
			roles: cat.roles,
			e2ee: cat.e2ee,
			channels: {}
		};

		categoryStore.addCategory(newCategory);
	}

	private handleChannelStart(value: StartChannel) {
		if (!value.channel) return;

		const chan = value.channel;
		const newChannel: Channel = {
			id: chan.id,
			server_id: chan.serverId,
			category_id: chan.categoryId,
			name: chan.name,
			description: chan.description,
			users: chan.users,
			roles: chan.roles,
			type: chan.type as ChannelTypes,
			position: chan.position,
			unread: false
		};

		channelStore.addChannel(newChannel);
	}

	private handleCategoryDelete(value: KillCategory) {
		const channels = channelStore.getCategoryChannels(value.serverId, value.categoryId);

		if (channels.find((chan) => chan.id === page.params.channel_id)) {
			const firstChan = channelStore.getFirstChannel(value.serverId);
			if (firstChan) goto(`/servers/${value.serverId}/channels/${firstChan}`);
		}

		categoryStore.deleteCategory(value.serverId, value.categoryId);
	}

	private handleChannelDelete(value: KillChannel) {
		if (!value.channel) return;

		const channel = value.channel;
		if (channel.id === page.params.channel_id) {
			const firstChan = channelStore.getFirstChannel(channel.serverId);
			if (firstChan) goto(`/servers/${channel.serverId}/channels/${firstChan}`);
		}

		channelStore.deleteChannel(channel.serverId, channel.categoryId, channel.id);
	}

	private handleRoleCreateOrEdit(value: CreateOrEditRole) {
		if (!value.role) return;

		const role = value.role;
		const newRole: Role = {
			id: role.id,
			members: [],
			name: role.name,
			color: role.color,
			abilities: role.abilities,
			position: role.position
		};

		const existingRole = serverStore.getRole(role.serverId, role.id);
		if (existingRole) {
			serverStore.editRole(role.serverId, newRole);
		} else {
			serverStore.addRole(role.serverId, newRole);
		}
	}

	private handleAddRoleMember(value: AddRoleMember) {
		if (!value.userId || !value.role) return;

		const member = serverStore.getMember(value.role.serverId, value.userId);
		if (member) member.roles.push(value.role.id);

		if (value.userId === userStore.user!.id) {
			serverStore.servers[value.role.serverId].user_roles.push(value.role.id);
		}
	}

	private handleRemoveRoleMember(value: RemoveRoleMember) {
		if (!value.userId || !value.role) return;

		const member = serverStore.getMember(value.role.serverId, value.userId);
		if (member) {
			member.roles = member.roles.filter((roleID) => roleID !== value.role.id);
		}

		if (value.userId === userStore.user!.id) {
			const server = serverStore.servers[value.role.serverId];
			server.user_roles = server.user_roles.filter((roleID) => roleID !== value.role.id);
		}
	}

	private handleRoleDelete(value: RemoveRole) {
		if (!value.role) return;
		serverStore.deleteRole(value.role.serverId, value.role.id);
	}

	private handleRoleMove(value: MoveRole) {
		const { movedRole, targetRole } = value;
		if (!movedRole?.id || !targetRole?.id || !movedRole?.serverId) return;

		const roles = serverStore.getRoles(movedRole.serverId);
		const moved = serverStore.getRole(movedRole.serverId, movedRole.id);
		const target = serverStore.getRole(movedRole.serverId, targetRole.id);

		if (moved && target) {
			moved.position = value.to;
			target.position = value.from;
			roles.sort((a, b) => a.position - b.position);
		}
	}

	private handleFriendRequest(value: FriendRequest) {
		if (!value.sender) return;

		const friend: Friend = {
			friendship_id: value.friendshipId,
			friendship_sender_id: value.sender.id,
			id: value.sender.id,
			display_name: value.sender.displayName,
			about_me: JSON.parse(new TextDecoder().decode(value.sender.aboutMe)),
			banner: value.sender.banner,
			avatar: value.sender.avatar,
			accepted: value.accepted,
			status: 'offline'
		};

		userStore.friends.push(friend);
	}

	private handleAcceptFriend(value: AcceptFriendRequest) {
		userStore.acceptFriend(value.friendshipId, value.channelId);
	}

	private handleRemoveFriend(value: RemoveFriend) {
		userStore.removeFriend({ friendshipID: value.friendshipId });
	}

	private handleAccountDeletion(value: AccountDeletion) {
		if (value.serverId !== '') {
			serverStore.deleteMember(value.serverId, value.userId);
		} else {
			userStore.removeFriend({ userID: value.userId });
		}
	}

	private handleServerAvatarChange(value: AvatarServerChange) {
		serverStore.updateAvatar(value.serverId, value.avatarUrl, value.bannerUrl);
	}

	private handleServerProfileChange(value: ProfileServerChange) {
		serverStore.updateProfile(value.serverId, {
			name: value.name,
			description: JSON.parse(new TextDecoder().decode(value.description)),
			public: value.public
		});
	}

	private handleChannelEdit(value: EditChannel) {
		if (!value.channel) return;

		const channel = value.channel;
		channelStore.editChannel(channel.id, {
			server_id: channel.serverId,
			name: channel.name,
			description: channel.description,
			users: channel.users,
			roles: channel.roles
		});

		if (channel.users?.length > 0 && !channel.users.includes(userStore.user!.id)) {
			goto(`/servers/${channel.serverId}`);
		}
	}

	private handleCategoryEdit(value: EditCategory) {
		if (!value.category) return;

		const category = value.category;
		categoryStore.editCategory(category.id, {
			server_id: category.serverId,
			name: category.name,
			users: category.users,
			roles: category.roles
		});
	}

	private handleServerDelete(value: KillServer) {
		if (!value.serverId) return;

		if (page.url.pathname.includes(value.serverId)) goto('/servers');
		serverStore.deleteServer(value.serverId);
	}

	private handleServerLeave(value: LeaveServer) {
		serverStore.deleteMember(value.serverId, value.userId);
	}

	private handleUserBan(value: BanUser) {
		serverStore.deleteMember(value.serverId, value.userId);

		if (value.userId === userStore.user?.id) {
			this.handleCurrentUserRestriction(value.serverId, value.reason, 'ban', "You've been banned");
		}
	}

	private handleUserKick(value: KickUser) {
		serverStore.deleteMember(value.serverId, value.userId);

		if (value.userId === userStore.user?.id) {
			this.handleCurrentUserRestriction(value.serverId, value.reason, 'kick', "You've been kicked");
		}
	}

	private handleCurrentUserRestriction(
		serverID: string,
		reason: string,
		restriction: string,
		title: string
	) {
		if (page.url.pathname.includes(serverID)) {
			goto('/servers');
			coreStore.restrictionDialog = {
				open: true,
				title,
				reason,
				restriction
			};
		}
		serverStore.deleteServer(serverID);
	}

	private parseAttachments(attachments: Uint8Array): any[] {
		const attachmentsStr = new TextDecoder().decode(attachments);
		return attachmentsStr.length > 0 ? JSON.parse(attachmentsStr) : [];
	}
}

export const ws = new WebsocketStore();
