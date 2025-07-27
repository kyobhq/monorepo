package types

import "encoding/json"

type ChatMessage struct {
	AuthorID  string          `json:"author_id" validate:"required"`
	ServerID  string          `json:"server_id" validate:"required"`
	ChannelID string          `json:"channel_id" validate:"required"`
	Content   json.RawMessage `json:"content" validate:"required"`
}
