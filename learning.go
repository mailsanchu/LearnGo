package main

import (
    "fmt"
    "time"
)

type ticker struct {
    period time.Duration
    ticker time.Ticker
}

type server struct {
    doneChan chan bool
    tickerA  ticker
    tickerB  ticker
}

func createTicker(period time.Duration) *ticker {
    return &ticker{period, *time.NewTicker(period)}
}

func (t *ticker) resetTicker() {
    t.ticker.Stop()
    t.ticker = *time.NewTicker(t.period)
}

func (s *server) listener() {
    start := time.Now()
    tickACount := 0
    fmt.Println("Elapsed: 0")
    for {
        select {
        case <-s.tickerA.ticker.C:
            elapsed := time.Since(start)
            fmt.Println("Elapsed: ", elapsed, " Ticker A")
            tickACount++
            if tickACount == 4 {
                s.doneChan <- true
            }
        case <-s.tickerB.ticker.C:
            s.tickerA.resetTicker()
            elapsed := time.Since(start)
            fmt.Println("Elapsed: ", elapsed, " Ticker B - Going to reset ticker A")
        }
    }
}

func main() {
    doneChan := make(chan bool)
    tickerA := createTicker(2 * time.Second)
    tickerB := createTicker(5 * time.Second)
    s := &server{doneChan, *tickerA, *tickerB}
    go s.listener()
    <-doneChan
}
