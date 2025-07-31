package types

type Ability string

const (
	ManageServer      Ability = "MANAGE_SERVER"
	ManageChannels    Ability = "MANAGE_CHANNELS"
	ManageRoles       Ability = "MANAGE_ROLES"
	ManageExpressions Ability = "MANAGE_EXPRESSIONS"
	ChangeNickname    Ability = "CHANGE_NICKNAME"
	ManageNicknames   Ability = "MANAGE_NICKNAMES"
	Ban               Ability = "BAN"
	Kick              Ability = "KICK"
	Mute              Ability = "MUTE"
	AttachFiles       Ability = "ATTACH_FILES"
	ManageMessages    Ability = "MANAGE_MESSAGES"
)

type CreateRoleParams struct {
	ServerID  string   `json:"server_id" validate:"required"`
	Position  int      `json:"position" validate:"min=0"`
	Name      string   `json:"name" validate:"required"`
	Color     string   `json:"color" validate:"required"`
	Abilities []string `json:"abilities" validate:"required"`
}
