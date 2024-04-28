package config

import "time"

type ApplicationConfig struct {
	Port     PortConfig     `json:"port"`
	Database DatabaseConfig `json:"database"`
	Rest     RestConfig     `json:"rest"`
	JWT      JWTConfig      `json:"jwt"`
}

type PortConfig struct {
	Service        int           `json:"service"`
	ServiceTimeout time.Duration `json:"servicetimeout"`
	BasePath       string        `json:"basepath"`
}

type DatabaseConfig struct {
	DBName  string        `json:"dbname"`
	Timeout time.Duration `json:"timeout"`
}

type RestConfig struct {
	Coincap CoincapConfig `json:"coincap"`
}

type CoincapConfig struct {
	BaseURL        string `json:"base_url"`
	AssetEndpoint  string `json:"asset_endpoint"`
	RatesEndpoint  string `json:"rates_endpoint"`
	TargetCurrency string `json:"target_currency"`
}

type JWTConfig struct {
	AccessTokenDuration  time.Duration `json:"access_token_duration"`
	RefreshTokenDuration time.Duration `json:"refresh_token_duration"`
	SecretKey            string        `json:"secret_key"`
}
