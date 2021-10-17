package calc

import (
	// "context"
	"fmt"
	"hw-async/domain"
	"math"
	"os"
	"sync"
	"time"
)

func FormCandleFromPrice(wg *sync.WaitGroup, prices <-chan domain.Price) <-chan domain.Candle {
	out := make(chan domain.Candle)
	var openTime time.Time
	candles := make(map[string]domain.Candle)
	go func() {
		defer func() {
			for _, c := range candles {
				out <- c
			}
			close(out)
			wg.Done()
		}()
		for price := range prices {
			time, _ := domain.PeriodTS(domain.CandlePeriod1m, price.TS)
			if openTime.Before(time) {
				openTime = time
				for _, c := range candles {
					out <- c
				}
				candles = make(map[string]domain.Candle)
			}
			if _, ok := candles[price.Ticker]; !ok {
				candles[price.Ticker] = domain.Candle{price.Ticker, domain.CandlePeriod1m, price.Value, price.Value,
					price.Value, price.Value, openTime}
			}
			candle := candles[price.Ticker]
			candle.High = math.Max(candle.High, price.Value)
			candle.Low = math.Min(candle.Low, price.Value)
			candle.Close = price.Value
			candles[price.Ticker] = candle
		}
	}()
	return out
}

func FormCandle(wg *sync.WaitGroup, candles <-chan domain.Candle, period domain.CandlePeriod) <-chan domain.Candle {
	out := make(chan domain.Candle)
	var openTime time.Time
	newCandles := make(map[string]domain.Candle)
	go func() {
		defer func() {
			for _, c := range newCandles {
				out <- c
			}
			close(out)
			wg.Done()
		}()
		for candle := range candles {
			time, _ := domain.PeriodTS(period, candle.TS)
			if openTime.Before(time) {
				openTime = time
				for _, c := range newCandles {
					out <- c
				}
				newCandles = make(map[string]domain.Candle)
			}
			if _, ok := newCandles[candle.Ticker]; !ok {
				newCandles[candle.Ticker] = domain.Candle{candle.Ticker, period, candle.Open, candle.High,
					candle.Low, candle.Close, openTime}
			}
			c := newCandles[candle.Ticker]
			c.High = math.Max(candle.High, c.High)
			c.Low = math.Min(candle.Low, c.Low)
			c.Close = candle.Close
			newCandles[candle.Ticker] = c
		}
	}()
	return out
}

func FormCandle2m(wg *sync.WaitGroup, candles <-chan domain.Candle) <-chan domain.Candle {
	return FormCandle(wg, candles, domain.CandlePeriod2m)
}

func FormCandle10m(wg *sync.WaitGroup, candles <-chan domain.Candle) <-chan domain.Candle {
	return FormCandle(wg, candles, domain.CandlePeriod10m)
}

func save(candle domain.Candle) {
	f, err := os.OpenFile(fmt.Sprintf("candles_%s.csv", candle.Period), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%v\n", candle))
}

func Save(wg *sync.WaitGroup, candles <-chan domain.Candle) <-chan domain.Candle {
	out := make(chan domain.Candle)
	go func() {
		defer wg.Done()
		for candle := range candles {
			out <- candle
			save(candle)
		}
		close(out)
	}()
	return out
}
