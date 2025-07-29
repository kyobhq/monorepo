package types

import (
	db "backend/db/gen_queries"
	"encoding/json"
)

type Crop struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Height int `json:"height"`
	Width  int `json:"width"`
}

type CreateServerParams struct {
	Name        string          `json:"name" validate:"required,min=1,max=20"`
	Description json.RawMessage `json:"description"`
	Public      bool            `json:"public"`
	Crop        Crop            `json:"crop" validate:"required"`
}

type JoinServerParams struct {
	InviteCode string `json:"invite_code" validate:"required"`
}

type UpdateServerProfileParams struct {
	Name        string          `json:"name" validate:"omitempty,min=1,max=20"`
	Description json.RawMessage `json:"description" validate:"omitempty"`
}

type UpdateServerAvatarParams struct {
	Avatar    string `json:"avatar" validate:"omitempty"`
	Banner    string `json:"banner" validate:"omitempty"`
	MainColor string `json:"main_color" validate:"omitempty"`
}

type ServerWithChannels struct {
	db.GetServersFromUserRow
	Channels map[string]db.Channel       `json:"channels"`
	Roles    []db.GetRolesFromServersRow `json:"roles"`
}
