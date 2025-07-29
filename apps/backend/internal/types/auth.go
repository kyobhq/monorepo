package types

type SignInParams struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}

type SignUpParams struct {
	Email       string `validate:"required,email" json:"email"`
	Username    string `validate:"required,max=20" json:"username"`
	DisplayName string `validate:"required,max=20" json:"display_name"`
	Password    string `validate:"required,min=8,max=254" json:"password"`
}
