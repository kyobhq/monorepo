package types

type CreateChannelParams struct {
	Position    int      `json:"position" validate:"min=0"`
	CategoryID  string   `json:"category_id" validate:"required"`
	ServerID    string   `json:"server_id" validate:"required"`
	Name        string   `json:"name" validate:"required,min=1,max=20"`
	Description string   `json:"description" validate:"omitempty"`
	Users       []string `json:"users" validate:"omitempty"`
	Roles       []string `json:"roles" validate:"omitempty"`
	E2EE        bool     `json:"e2ee"`
	Type        string   `json:"type" validate:"required,oneof=textual voice gallery kanban"`
}

type CreateCategoryParams struct {
	ServerID string   `json:"server_id" validate:"required"`
	Name     string   `json:"name" validate:"required,min=1,max=20"`
	Position int      `json:"position" validate:"min=0"`
	Users    []string `json:"users" validate:"omitempty"`
	Roles    []string `json:"roles" validate:"omitempty"`
	E2EE     bool     `json:"e2ee"`
}

type EditChannelParams struct {
	ServerID    string   `json:"server_id"`
	Name        string   `json:"name" validate:"required,min=1,max=20"`
	Description string   `json:"description" validate:"omitempty"`
	Users       []string `json:"users" validate:"omitempty"`
	Roles       []string `json:"roles" validate:"omitempty"`
	Private     bool     `json:"private" validate:"omitempty"`
}

type EditCategoryParams struct {
	ServerID string   `json:"server_id"`
	Name     string   `json:"name" validate:"required,min=1,max=20"`
	Users    []string `json:"users" validate:"omitempty"`
	Roles    []string `json:"roles" validate:"omitempty"`
}

type PinChannelParams struct {
	ServerID string `json:"server_id" validate:"required"`
	Position int    `json:"position" validate:"min=0"`
}

type DeleteChannelParams struct {
	ServerID   string `json:"server_id" validate:"required"`
	CategoryID string `json:"category_id" validate:"required"`
}

type DeleteCategoryParams struct {
	ServerID    string   `json:"server_id" validate:"required"`
	ChannelsIDs []string `json:"channels_ids" validate:"omitempty"`
}
