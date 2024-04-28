package model

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserToken struct {
	UserId         int    `json:"userId"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	ExpirationTime int64  `json:"expiration_time"`
}

type UserAsset struct {
	ID      int    `json:"id"`
	UserId  int    `json:"userId"`
	AssetId string `json:"assetId"`
}

type Asset struct {
	AssetId string  `json:"assetId"`
	Price   float64 `json:"price"`
}
