package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
    // "encoding/json"
	// "time"
    "bot/repository/exchange"
)


func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

    k := exchange.KrakenService{}
    err := k.Connect("wss://futures.kraken.com/ws/v1")
    out, err := k.Subscribe([]string{"PI_XBTUSD"})
    if err != nil {
        log.Println(err)
    }
    go func() {
		for msg := range(out) {
            log.Println(msg)
        }
    }()
    <-interrupt
    k.Close()

}
