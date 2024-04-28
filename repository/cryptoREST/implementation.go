package cryptoREST

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/michaelwongycn/crypto-tracker/domain/model"
	"github.com/michaelwongycn/crypto-tracker/domain/response"
	"github.com/michaelwongycn/crypto-tracker/lib/log"
)

const (
	apiRequestFailedErrorMsg   = "API request failed with status code"
	invalidAPIResponseErrorMsg = "error when decoding response json"
	errorAccessingAPIErrorMsg  = "error when accessing external API"
	errorParsingPriceErrorMsg  = "error when parsing price string"
)

type cryptoRESTImpl struct {
	timeout        time.Duration
	baseURL        string
	assetEndpoint  string
	ratesEndpoint  string
	targetCurrency string
}

func NewCryptoRESTImpl(timeout time.Duration, baseURL, assetEndpoint, ratesEndpoint, targetCurrency string) CryptoRESTInterface {
	return &cryptoRESTImpl{
		timeout:        timeout * time.Second,
		baseURL:        baseURL,
		assetEndpoint:  assetEndpoint,
		ratesEndpoint:  ratesEndpoint,
		targetCurrency: targetCurrency,
	}
}

func (r *cryptoRESTImpl) IsValidAsset(ctx context.Context, asset string) (bool, error) {
	resp, err := http.Get(r.baseURL + r.assetEndpoint + asset)
	if err != nil {
		log.PrintLogErr(ctx, errorAccessingAPIErrorMsg, err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.PrintLogAPIErr(ctx, apiRequestFailedErrorMsg, resp.StatusCode)
		return false, fmt.Errorf("%s: %d", apiRequestFailedErrorMsg, resp.StatusCode)
	}

	var APIResponse struct {
		Data response.AssetValidationDataResponse `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&APIResponse)
	if err != nil {
		log.PrintLogErr(ctx, invalidAPIResponseErrorMsg, err)
		return false, err
	}

	return APIResponse.Data.ID == asset, nil
}

func (r *cryptoRESTImpl) GetAssetsPrice(ctx context.Context, userAssets *[]model.UserAsset) (*[]model.Asset, error) {
	resp, err := http.Get(r.baseURL + r.ratesEndpoint + r.targetCurrency)
	if err != nil {
		log.PrintLogErr(ctx, errorAccessingAPIErrorMsg, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.PrintLogAPIErr(ctx, apiRequestFailedErrorMsg, resp.StatusCode)
		return nil, fmt.Errorf("%s: %d", apiRequestFailedErrorMsg, resp.StatusCode)
	}

	var APIResponse struct {
		Data response.CurrencyRateDataResponse `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&APIResponse)
	if err != nil {
		log.PrintLogErr(ctx, invalidAPIResponseErrorMsg, err)
		return nil, err
	}

	rateIDR, err := strconv.ParseFloat(APIResponse.Data.RateUSD, 64)
	if err != nil {
		log.PrintLogErr(ctx, errorParsingPriceErrorMsg, err)
		return nil, err
	}

	var data []model.Asset
	for _, userAsset := range *userAssets {
		var asset model.Asset
		resp, err := http.Get(r.baseURL + r.assetEndpoint + userAsset.AssetId)
		if err != nil {
			log.PrintLogErr(ctx, errorAccessingAPIErrorMsg, err)
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.PrintLogAPIErr(ctx, apiRequestFailedErrorMsg, resp.StatusCode)
			return nil, fmt.Errorf("%s: %d", apiRequestFailedErrorMsg, resp.StatusCode)
		}

		var APIResponse struct {
			Data response.AssetValidationDataResponse `json:"data"`
		}

		err = json.NewDecoder(resp.Body).Decode(&APIResponse)
		if err != nil {
			log.PrintLogErr(ctx, invalidAPIResponseErrorMsg, err)
			return nil, err
		}

		price, err := strconv.ParseFloat(APIResponse.Data.PriceUSD, 64)
		if err != nil {
			log.PrintLogErr(ctx, errorParsingPriceErrorMsg, err)
			return nil, err
		}

		asset.AssetId = userAsset.AssetId
		asset.Price = price / rateIDR

		data = append(data, asset)
	}

	return &data, nil
}
