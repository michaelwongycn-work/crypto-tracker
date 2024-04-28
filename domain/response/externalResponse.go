package response

type AssetValidationDataResponse struct {
	ID                string `json:"id"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	PriceUSD          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
}

type CurrencyRateDataResponse struct {
	ID             string `json:"id"`
	Symbol         string `json:"symbol"`
	CurrencySymbol string `json:"currencySymbol"`
	Type           string `json:"type"`
	RateUSD        string `json:"rateUsd"`
}
