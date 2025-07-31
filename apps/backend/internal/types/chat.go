package types

import "encoding/json"

type CreateMessageParams struct {
	ServerID         string          `json:"server_id" validate:"required"`
	ChannelID        string          `json:"channel_id" validate:"required"`
	Content          json.RawMessage `json:"content" validate:"required"`
	Everyone         bool            `json:"everyone"`
	MentionsUsers    []string        `json:"mentions_users"`
	MentionsRoles    []string        `json:"mentions_roles"`
	MentionsChannels []string        `json:"mentions_channels"`
	Attachments      json.RawMessage `json:"attachments"`
}

type File struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Filename string `json:"file_name"`
	Filesize string `json:"file_size"`
	Type     string `json:"type"`
}
