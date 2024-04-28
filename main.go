package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/michaelwongycn/crypto-tracker/controller"
	"github.com/michaelwongycn/crypto-tracker/handler"
	"github.com/michaelwongycn/crypto-tracker/lib/auth"
	"github.com/michaelwongycn/crypto-tracker/lib/cache"
	"github.com/michaelwongycn/crypto-tracker/lib/cfg"
	"github.com/michaelwongycn/crypto-tracker/lib/db"
	"github.com/michaelwongycn/crypto-tracker/repository/cryptoDB"
	"github.com/michaelwongycn/crypto-tracker/repository/cryptoREST"
	"github.com/michaelwongycn/crypto-tracker/usecase/user"
)

func main() {
	cfg, err := cfg.ReadConfig()
	if err != nil {
		log.Printf("Error reading config: %v\n", err)
	}

	auth.SetAuthConfig(cfg.JWT.SecretKey, cfg.JWT.AccessTokenDuration, cfg.JWT.RefreshTokenDuration)
	cache.InitializeNewCache(*cfg)
	db, err := db.Connect(cfg.Database.Timeout, cfg.Database.DBName)
	if err != nil {
		log.Printf("Error connecting to DB: %v\n", err)
	}

	cryptoDB := cryptoDB.NewCryptoDBImpl(60, db)
	cryptoREST := cryptoREST.NewCryptoRESTImpl(60, cfg.Rest.Coincap.BaseURL, cfg.Rest.Coincap.AssetEndpoint, cfg.Rest.Coincap.RatesEndpoint, cfg.Rest.Coincap.TargetCurrency)

	userUsecase := user.NewUserImpl(cryptoDB, cryptoREST, cfg.JWT.RefreshTokenDuration)

	controller := controller.NewControllerImpl(userUsecase)

	handler := handler.NewHandler(60, controller)

	rest := handler.StartRoute()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Printf("Shutdown Application ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cancel()
		db.Close()
	}()

	if err := rest.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v", err)
	}
	log.Printf("Application Stopped")
}
