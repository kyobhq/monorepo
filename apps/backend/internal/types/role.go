package types

type Ability string

const (
	ViewChannels      Ability = "VIEW_CHANNELS"
	ManageChannels    Ability = "MANAGE_CHANNELS"
	ManageRoles       Ability = "MANAGE_ROLES"
	ManageServer      Ability = "MANAGE_SERVER"
	CreateInvite      Ability = "CREATE_INVITE"
	ChangeNickname    Ability = "CHANGE_NICKNAME"
	ManageNicknames   Ability = "MANAGE_NICKNAMES"
	TimeoutMembers    Ability = "TIMEOUT_MEMBERS"
	KickMembers       Ability = "KICK_MEMBERS"
	BanMembers        Ability = "BAN_MEMBERS"
	SendMessages      Ability = "SEND_MESSAGES"
	AttachFiles       Ability = "ATTACH_FILES"
	AddReactions      Ability = "ADD_REACTIONS"
	UsePersonalEmojis Ability = "USE_PERSONAL_EMOJIS"
	MentionEveryone   Ability = "MENTION_EVERYONE"
	ManageMessages    Ability = "MANAGE_MESSAGES"
	Connect           Ability = "CONNECT"
	Speak             Ability = "SPEAK"
	Video             Ability = "VIDEO"
	MuteMembers       Ability = "MUTE_MEMBERS"
	DeafenMembers     Ability = "DEAFEN_MEMBERS"
	MoveMembers       Ability = "MOVE_MEMBERS"
	Administrator     Ability = "ADMINISTRATOR"
)

type CreateRoleParams struct {
	RoleID    string   `json:"id" validate:"required"`
	ServerID  string   `json:"server_id" validate:"required"`
	Position  int      `json:"position" validate:"min=0"`
	Name      string   `json:"name" validate:"required"`
	Color     string   `json:"color" validate:"required"`
	Abilities []string `json:"abilities" validate:"required"`
}

type DeleteRoleParams struct {
	ServerID string `json:"server_id" validate:"required"`
	RoleID   string `json:"role_id" validate:"required"`
}

type ChangeRoleMemberParams struct {
	ServerID string `json:"server_id" validate:"required"`
	RoleID   string `json:"role_id" validate:"required"`
	UserID   string `json:"user_id" validate:"omitempty"`
}

type MoveRoleMemberParams struct {
	ServerID     string `json:"server_id" validate:"required"`
	TargetRoleID string `json:"target_role_id" validate:"required"`
	MovedRoleID  string `json:"moved_role_id" validate:"required"`
	From         int    `json:"from" validate:"omitempty"`
	To           int    `json:"to" validate:"omitempty"`
}
