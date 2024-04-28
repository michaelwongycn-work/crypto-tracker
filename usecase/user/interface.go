package user

import (
	"context"

	"github.com/michaelwongycn/crypto-tracker/domain/model"
)

type UserUsecase interface {
	Login(ctx context.Context, email, password string) (*string, *string, error)
	Register(ctx context.Context, email, password string) error
	Logout(ctx context.Context, accessToken string, userId int) error
	RefreshToken(ctx context.Context, refreshToken string, userId int) (*string, *string, error)
	GetUserAssetsByUserId(ctx context.Context, userId int) (*[]model.Asset, error)
	InsertUserAsset(ctx context.Context, userId int, assetId string) error
	DeleteUserAsset(ctx context.Context, userId int, assetId string) error
}
