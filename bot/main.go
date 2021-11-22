package main

import (
	"bot/repository/krakenapi"
	"flag"
	"log"
	"os"
	"os/signal"
)

//type config struct {
//	PublicKey string `yaml:"publickey"`
//	PrivateKey string `yaml:"privatekey"`
//	Timeout time.Duration`yaml:"duration"`
//}

func main() {

	var filePathFlag = flag.String("config", "", "path to file")
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	k := krakenapi.NewFromConfig(*filePathFlag)

	// TEST WEBSOCKET
	//ctx, cancel := context.WithCancel(context.Background())
	//err = k.Connect(3)
	//if err != nil {
	//	log.Println(err)
	//}
	//out, err := k.Subscribe("PI_ETHUSD")
	//if err != nil {
	//	log.Println(err)
	//}
	//go func() {
	//	for msg := range out {
	//		log.Println(msg)
	//	}
	//}()
	//<-interrupt
	//k.Stop()
	//time.Sleep(2 * time.Second)

	//TEST REST
	//resp, err := k.SendOrder("PI_ETHUSD", 4411, 10, krakenapi.Sell)
	//if err != nil {
	//	log.Println(err)
	//}
	//defer resp.Body.Close()
	//bodyBytes, err := io.ReadAll(resp.Body)
	//log.Println(string(bodyBytes))

	cancelResp, _ := k.CancelOrders()
	log.Println(cancelResp.CancelStatus.Status)
	log.Println(cancelResp.Result)

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
