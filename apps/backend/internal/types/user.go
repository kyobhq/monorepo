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
}

type UpdatePasswordParams struct {
	Current string `json:"current" validate:"required,min=8,max=254"`
	New     string `json:"new" validate:"required,min=8,max=254"`
	Confirm string `json:"confirm" validate:"required,min=8,max=254"`
}

type Link struct {
	Label string `json:"label" validate:"omitempty,min=1,max=20"`
	URL   string `json:"url" validate:"omitempty,url"`
}

type Fact struct {
	Label string `json:"label" validate:"omitempty,min=1,max=20"`
	Value string `json:"value" validate:"omitempty,min=1,max=20"`
}
