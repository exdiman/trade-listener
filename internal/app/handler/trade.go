package handler

import (
	"log"
	"strconv"

	"github.com/exdiman/trade-listener/internal/app/binance"
	"github.com/exdiman/trade-listener/internal/app/model"
	"github.com/exdiman/trade-listener/internal/app/poloniex"
)

type BinanceTrade struct {
	Orig *binance.Trade
}

type PoloniexTrade struct {
	Orig *poloniex.Trade
}

type TradeContract interface {
	GetBaseAsset() string
	GetQuoteAsset() string
	GetTradeId() string
	GetSide() model.Side
	GetPrice() float64
	GetAmount() float64
	GetTimestampMs() int
}

func MakeTrade(t TradeContract, exchange model.Exchange) *model.Trade {
	return &model.Trade{
		Exchange:    exchange,
		BaseAsset:   t.GetBaseAsset(),
		QuoteAsset:  t.GetQuoteAsset(),
		TradeID:     t.GetTradeId(),
		Side:        t.GetSide(),
		Price:       t.GetPrice(),
		Amount:      t.GetAmount(),
		TimestampMs: t.GetTimestampMs(),
	}
}

/** BinanceTrade getters */

func (t BinanceTrade) GetBaseAsset() string {
	return t.Orig.BaseAsset
}

func (t BinanceTrade) GetQuoteAsset() string {
	return t.Orig.QuoteAsset
}

func (t BinanceTrade) GetTradeId() string {
	return strconv.Itoa(t.Orig.TradeID)
}

func (t BinanceTrade) GetSide() model.Side {
	return model.SideBuy
}

func (t BinanceTrade) GetPrice() float64 {
	return t.Orig.Price
}

func (t BinanceTrade) GetAmount() float64 {
	return t.Orig.Quantity
}

func (t BinanceTrade) GetTimestampMs() int {
	return t.Orig.TradeTime
}

/** PoloniexTrade getters */

func (t PoloniexTrade) GetBaseAsset() string {
	return t.Orig.BaseAsset
}

func (t PoloniexTrade) GetQuoteAsset() string {
	return t.Orig.QuoteAsset
}

func (t PoloniexTrade) GetTradeId() string {
	return t.Orig.TradeID
}

func (t PoloniexTrade) GetSide() model.Side {
	if t.Orig.Side == 0.0 {
		return model.SideSell
	}
	return model.SideBuy
}

func (t PoloniexTrade) GetPrice() float64 {
	f, err := strconv.ParseFloat(t.Orig.Price, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func (t PoloniexTrade) GetAmount() float64 {
	f, err := strconv.ParseFloat(t.Orig.Size, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func (t PoloniexTrade) GetTimestampMs() int {
	return int(t.Orig.Timestamp) * 1000
}
