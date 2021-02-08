package main

import (
	"context"
	"fmt"
	"github.com/exdiman/trade-listener/internal/app/model"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/exdiman/trade-listener/internal/app/binance"
	"github.com/exdiman/trade-listener/internal/app/handler"
	"github.com/exdiman/trade-listener/internal/app/poloniex"
	"github.com/exdiman/trade-listener/internal/app/sqlstore"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())

	// TODO работа с конфигами
	db, err := sqlx.Connect("mysql", "gouser:password@(localhost:33061)/trade_listener")
	if err != nil {
		log.Fatalln(err)
	}

	store := sqlstore.New(db)
	store.PrepareDb() // TODO заменить на миграции

	// канал для сбора всех сделок
	tradeCh := make(chan *model.Trade, 100)
	defer close(tradeCh)

	// обработчик поступающих сделок
	done := handler.StoreTrades(ctx, tradeCh, store)

	var wg sync.WaitGroup

	// получение сделок на бирже Binance
	binanceClient := binance.NewClient(ctx)
	binanceTradeCh, binanceDone := binanceClient.ListenTrades("BTC", "USDT")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case trade := <-binanceTradeCh:
				tradeCh <- handler.MakeTrade(handler.BinanceTrade{Orig: trade}, "Binance")
			case <-binanceDone:
				return
			}
		}
	}()

	// получение сделок на бирже Poloniex
	poloniexClient := poloniex.NewClient(ctx)
	poloniexTradeCh, poloniexDone := poloniexClient.ListenTrades("USDT", "BTC")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case trade := <-poloniexTradeCh:
				tradeCh <- handler.MakeTrade(handler.PoloniexTrade{Orig: trade}, "Poloniex")
			case <-poloniexDone:
				return
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGKILL)

	<-interrupt
	fmt.Println("canceling context")
	cancelCtx()
	fmt.Println("waiting for canceling")
	wg.Wait()
	<-done
	fmt.Println("exit")

	return
}
