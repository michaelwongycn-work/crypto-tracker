package cryptoREST

import (
	"context"

	"github.com/michaelwongycn/crypto-tracker/domain/model"
)

type CryptoRESTInterface interface {
	IsValidAsset(ctx context.Context, asset string) (bool, error)
	GetAssetsPrice(ctx context.Context, userAssets *[]model.UserAsset) (*[]model.Asset, error)
}
