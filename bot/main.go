package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	// "encoding/json"
	// "bot/repository/exchange"
    // "context"
	"bot/repository/telegram"
	"time"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)


	// ctx, cancel := context.WithCancel(context.Background())
	// k := exchange.KrakenService{}
	// err := k.Connect(ctx, "wss://futures.kraken.com/ws/v1")
	// out, err := k.Subscribe([]string{"PI_XBTUSD"})
	// if err != nil {
	// 	log.Println(err)
	// }
	// go func() {
	// 	for msg := range out {
	// 		log.Println(msg)
	// 	}
	// }()
	// <-interrupt
	// k.Close()
	// cancel()
	// time.Sleep(2 * time.Second)

	tgBot, err := telegram.NewBot("ef")
    if err != nil {
        log.Println(err)
        return
    }
	tgBot.Start()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			tgBot.Notify("refer")
		}
	}()
	<-interrupt
	tgBot.Stop()
	time.Sleep(1 * time.Second)
}
