package main

import (
	"bot/repository/krakenapi"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//ctx, cancel := context.WithCancel(context.Background())
	k := krakenapi.New("kDZBu4RF0DeXPNtM9jYdXfJtevCL+JdCb01+5AIKwm/oCc6ipxlXY8Zd",
		"/1COOB8P/vcQhQSs3048PL/oasEexmYOh/znNttPKY7PjiYsReS2/ioI5kgDTyVG99ARc6ROMreqLhHELx93LcDE",
		10*time.Second)

	// TEST WEBSOCKET
	//err := k.Connect(ctx, "wss://futures.kraken.com/ws/v1")
	//out, err := k.Subscribe([]string{"PI_XBTUSD"})
	//if err != nil {
	//	log.Println(err)
	//}
	//go func() {
	//	for msg := range out {
	//		log.Println(msg)
	//	}
	//}()
	//<-interrupt
	//k.Close()
	//cancel()
	//time.Sleep(2 * time.Second)

	//TEST REST
	resp, err := k.CancelOrders()
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	log.Println(string(bodyBytes))

	//TEST TELEGRAM
	//tgBot, err := telegram.NewBot("ef")
	//if err != nil {
	//    log.Println(err)
	//    return
	//}
	//tgBot.Start()
	//
	//go func() {
	//	for {
	//		time.Sleep(5 * time.Second)
	//		tgBot.Notify("refer")
	//	}
	//}()
	//<-interrupt
	//tgBot.Stop()
	//time.Sleep(1 * time.Second)
}
