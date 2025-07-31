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
	Position    int             `json:"position" validate:"min=0"`
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

type ServerWithCategories struct {
	db.GetServersFromUserRow
	Categories map[string]CategoryWithChannels `json:"categories"`
	Roles      []db.GetRolesFromServersRow     `json:"roles"`
}

type CategoryWithChannels struct {
	db.ChannelCategory
	Channels map[string]db.Channel `json:"channels"`
}
