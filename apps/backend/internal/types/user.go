package types

import "encoding/json"

type UpdateAccountParams struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,min=1,max=20"`
}

type UpdateAvatarParams struct {
	Avatar    string `json:"avatar" validate:"omitempty"`
	Banner    string `json:"banner" validate:"omitempty"`
	MainColor string `json:"main_color" validate:"omitempty"`
}

type UpdateProfileParams struct {
	DisplayName string          `json:"display_name" validate:"omitempty,min=1,max=20"`
	About       json.RawMessage `json:"about" validate:"omitempty"`
	Facts       json.RawMessage `json:"facts" validate:"omitempty"`
	Links       json.RawMessage `json:"links" validate:"omitempty"`
}

type Setup struct {
	Servers map[string]ServerWithCategories `json:"servers"`
}
