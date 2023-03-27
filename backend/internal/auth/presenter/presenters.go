package presenter

type UserSignIn struct {
	Email    string `json:"email" validate:"required" example:"hiennguyen9874@gmail.com"`
	Password string `json:"password" validate:"required,min=8" example:"password"`
}

type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

type PublicKey struct {
	PublicKeyAccessToken  string `json:"public_key_access_token,omitempty"`
	PublicKeyRefreshToken string `json:"public_key_refresh_token,omitempty"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required" example:"hiennguyen9874@gmail.com"`
}

type ResetPassword struct {
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8" example:"password"`
}
