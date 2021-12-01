package main

import (
	"bot/handlers"
	"bot/krakenapi"
	"bot/modelapi"
	"bot/repository"
	"bot/service"
	"bot/telegramapi"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	//logger := log.New()
	//logger.SetLevel(log.DebugLevel)

	// repository
	dsn := "postgres://user:passwd@localhost:5432/orders" +
		"?sslmode=disable"

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Println(err)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	repo := repository.New(pool)

	// notifier
	telegramNotifier, _ := telegramapi.NewWithCreds("tg_creds")

	// kraken api
	krakenAPI, _ := krakenapi.NewWithConfig("kraken_config.yaml")

	// model service
	url := "http://localhost:7070/v1/models/trade_model:predict"
	modelService := modelapi.New(url)

	// service
	tradebot := service.New(krakenAPI, telegramNotifier, repo, modelService,
		"PI_XBTUSD", 10, 1, 0.5, 20, 1)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	tradebotHandler := handlers.New(tradebot)
	r.Mount("/", tradebotHandler.Routes())

	http.ListenAndServe(":3000", r)
}
