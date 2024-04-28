package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/michaelwongycn/crypto-tracker/domain/model"
	"github.com/michaelwongycn/crypto-tracker/lib/auth"
	"github.com/michaelwongycn/crypto-tracker/lib/cache"
	"github.com/michaelwongycn/crypto-tracker/repository/cryptoDB"
	"github.com/michaelwongycn/crypto-tracker/repository/cryptoREST"
)

type userImpl struct {
	dbCrypto             cryptoDB.CryptoDBInterface
	restCrypto           cryptoREST.CryptoRESTInterface
	refreshTokenDuration time.Duration
}

func NewUserImpl(dbCrypto cryptoDB.CryptoDBInterface, restCrypto cryptoREST.CryptoRESTInterface, refreshTokenDuration time.Duration) UserUsecase {
	return &userImpl{
		dbCrypto:             dbCrypto,
		restCrypto:           restCrypto,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (u *userImpl) Login(ctx context.Context, email, password string) (*string, *string, error) {
	currTime := time.Now()
	// TODO: encrypt Password
	user, err := u.dbCrypto.GetUserByEmailAndPassword(ctx, email, password)
	if err != nil {
		return nil, nil, err
	}

	oldUserToken, err := u.dbCrypto.GetUserToken(ctx, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}

	accessToken, refreshToken, err := auth.CreateToken(currTime, user.ID)
	if err != nil {
		return nil, nil, err
	}

	err = u.dbCrypto.InsertUserToken(ctx, user.ID, accessToken, refreshToken, currTime.Add(time.Minute*u.refreshTokenDuration).Unix())
	if err != nil {
		return nil, nil, err
	}

	if oldUserToken != nil {
		cache.DeleteCache(oldUserToken.AccessToken)
	}

	cache.SetCache(accessToken, refreshToken)
	return &accessToken, &refreshToken, nil
}

func (u *userImpl) Register(ctx context.Context, email, password string) error {
	// TODO: encrypt Password
	return u.dbCrypto.InsertUser(ctx, email, password)
}

func (u *userImpl) Logout(ctx context.Context, accessToken string, userId int) error {
	cache.DeleteCache(accessToken)
	return u.dbCrypto.DeleteUserToken(ctx, userId)
}

func (u *userImpl) RefreshToken(ctx context.Context, refreshToken string, userId int) (*string, *string, error) {
	userToken, err := u.dbCrypto.GetUserToken(ctx, userId)
	if err != nil {
		return nil, nil, err
	}

	currTime := time.Now()

	if refreshToken != userToken.RefreshToken || userToken.ExpirationTime < currTime.Unix() {
		return nil, nil, errors.New("invalid refresh token")
	}

	newAccessToken, newRefreshToken, err := auth.CreateToken(currTime, userId)
	if err != nil {
		return nil, nil, err
	}

	err = u.dbCrypto.InsertUserToken(ctx, userId, newAccessToken, newRefreshToken, currTime.Add(time.Minute*u.refreshTokenDuration).Unix())
	if err != nil {
		return nil, nil, err
	}

	cache.DeleteCache(userToken.AccessToken)
	cache.SetCache(newAccessToken, newRefreshToken)
	return &newAccessToken, &newRefreshToken, nil
}

func (u *userImpl) GetUserAssetsByUserId(ctx context.Context, userId int) (*[]model.Asset, error) {
	userAssets, err := u.dbCrypto.GetUserAssetsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	assetsPrice, err := u.restCrypto.GetAssetsPrice(ctx, userAssets)
	if err != nil {
		return nil, err
	}

	return assetsPrice, nil
}

func (u *userImpl) InsertUserAsset(ctx context.Context, userId int, assetId string) error {
	_, err := u.restCrypto.IsValidAsset(ctx, assetId)
	if err != nil {
		return err
	}
	return u.dbCrypto.InsertUserAsset(ctx, userId, assetId)
}

func (u *userImpl) DeleteUserAsset(ctx context.Context, userId int, assetId string) error {
	_, err := u.restCrypto.IsValidAsset(ctx, assetId)
	if err != nil {
		return err
	}
	return u.dbCrypto.DeleteUserAsset(ctx, userId, assetId)
}
