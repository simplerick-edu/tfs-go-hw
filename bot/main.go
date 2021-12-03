package main

import (
	"flag"
	"net/http"
	"os"
	// not-std
	"bot/handlers"
	"bot/krakenapi"
	"bot/modelapi"
	"bot/repository"
	"bot/service"
	"bot/telegramapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var dsn string
var telegramToken string
var modelServiceURL string
var krakenapiConfig krakenapi.Config
var botConfig service.Parameters

func init() {
	dsnPath := flag.String("dsn_path", "", "path to dsn")
	krakenConfigPath := flag.String("kraken_config_path", "", "path to yaml file with kraken config")
	telegramCredsPath := flag.String("telegram_creds_path", "", "path to file with telegram bot token")
	modelConfigPath := flag.String("model_config_path", "", "path to file with model service url")
	botConfigPath := flag.String("bot_config_path", "", "path to yaml file with bot parameters")
	flag.Parse()
	data, err := os.ReadFile(*dsnPath)
	dsn = string(data)
	data, err = os.ReadFile(*telegramCredsPath)
	telegramToken = string(data)
	data, err = os.ReadFile(*modelConfigPath)
	modelServiceURL = string(data)
	data, err = os.ReadFile(*botConfigPath)
	err = yaml.Unmarshal(data, &botConfig)
	if err != nil {
		log.Fatal(err)
	}
	data, err = os.ReadFile(*krakenConfigPath)
	err = yaml.Unmarshal(data, &krakenapiConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// repository
	pool, err := repository.NewPool(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	repo := repository.New(pool)

	// notifier
	telegramNotifier, err := telegramapi.NewWithCreds(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	// model service
	modelService := modelapi.New(modelServiceURL)

	// kraken api
	krakenAPI := krakenapi.NewWithConfig(krakenapiConfig)

	// service
	tradeBot := service.New(krakenAPI, telegramNotifier, repo, modelService, botConfig)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	tradebotHandler := handlers.New(tradeBot)
	r.Mount("/", tradebotHandler.Routes())
	log.Fatal(http.ListenAndServe(":3000", r))
}
