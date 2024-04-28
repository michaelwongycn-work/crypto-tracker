package cryptoDB

import (
	"context"

	"github.com/michaelwongycn/crypto-tracker/domain/model"
)

type CryptoDBInterface interface {
	GetUserByEmailAndPassword(ctx context.Context, email, password string) (*model.User, error)
	InsertUser(ctx context.Context, email, password string) error
	GetUserToken(ctx context.Context, userId int) (*model.UserToken, error)
	InsertUserToken(ctx context.Context, userId int, accessToken, refreshToken string, expirationTime int64) error
	DeleteUserToken(ctx context.Context, userId int) error

	GetUserAssetsByUserId(ctx context.Context, userId int) (*[]model.UserAsset, error)
	InsertUserAsset(ctx context.Context, userId int, assetId string) error
	DeleteUserAsset(ctx context.Context, userId int, assetId string) error
}
