package types

import (
	db "backend/db/gen_queries"
	"encoding/json"
)

type UpdateEmailParams struct {
	Email string `json:"email" validate:"omitempty,email"`
}

type UpdateAvatarParams struct {
	CropAvatar Crop   `json:"crop_avatar" validate:"omitempty"`
	CropBanner Crop   `json:"crop_banner" validate:"omitempty"`
	MainColor  string `json:"main_color" validate:"omitempty"`
}

type UpdateProfileParams struct {
	Username    string          `json:"username" validate:"omitempty,min=1,max=20"`
	DisplayName string          `json:"display_name" validate:"omitempty,min=1,max=20"`
	About       json.RawMessage `json:"about_me" validate:"omitempty"`
	Facts       json.RawMessage `json:"facts" validate:"omitempty"`
	Links       json.RawMessage `json:"links" validate:"omitempty"`
}

type Setup struct {
	User    *db.User                        `json:"user"`
	Servers map[string]ServerWithCategories `json:"servers"`
	Friends []db.GetFriendsRow              `json:"friends"`
	Emojis  []db.GetEmojisRow               `json:"emojis"`
}

type UpdatePasswordParams struct {
	Current string `json:"current" validate:"required,min=8,max=254"`
	New     string `json:"new" validate:"required,min=8,max=254"`
	Confirm string `json:"confirm" validate:"required,min=8,max=254"`
}

type GetUserProfileParams struct {
	ServerID string `json:"server_id" validate:"omitempty"`
}

type Link struct {
	Label string `json:"label" validate:"omitempty,min=1,max=20"`
	URL   string `json:"url" validate:"omitempty,url"`
}

type Fact struct {
	Label string `json:"label" validate:"omitempty,min=1,max=20"`
	Value string `json:"value" validate:"omitempty,min=1,max=20"`
}

type UploadEmojiParams struct {
	Shortcodes []string `validate:"required,max=20,dive,emoji_shortcode" json:"shortcodes"`
}

type UpdateEmojiParams struct {
	Shortcode string `validate:"required,max=20,emoji_shortcode" json:"shortcode"`
}

type EmojiResponse struct {
	ID        string `json:"id"`
	Shortcode string `json:"shortcode"`
	URL       string `json:"url"`
}

type SyncParams struct {
	ChannelIDs     []string          `json:"channel_ids"`
	LastMessageIDs []string          `json:"last_message_ids"`
	MentionIDs     []json.RawMessage `json:"mention_ids"`
}
