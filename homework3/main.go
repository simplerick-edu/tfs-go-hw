package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"hw-async/calc"
	"hw-async/generator"
	"os"
	"os/signal"
	"sync"
	"time"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	wg := sync.WaitGroup{}

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	wg.Add(6)
	prices := pg.Prices(ctx)
	candles := calc.FormCandleFromPrice(&wg, prices)
	candles = calc.Save(&wg, candles)
	candles = calc.FormCandle2m(&wg, candles)
	candles = calc.Save(&wg, candles)
	candles = calc.FormCandle10m(&wg, candles)
	candles = calc.Save(&wg, candles)

CheckStop:
	for {
		select {
		case <-stop:
			cancel()
		case _, ok := <-candles:
			if !ok {
				break CheckStop
			}
		}
	}
	logger.Info("all goroutines terminated")
	wg.Wait()
	logger.Info("main process terminated")
}
