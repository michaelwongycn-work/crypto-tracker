package request

type UserRegisterRequest struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type UserAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserInsertAssetRequest struct {
	AssetID string `json:"assetId"`
}
