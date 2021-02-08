package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Trade struct {
	BaseAsset    string
	QuoteAsset   string
	Type         string  `json:"e"`        // Event type
	EventTime    int     `json:"E"`        // Event time
	Symbol       string  `json:"s"`        // Symbol
	TradeID      int     `json:"a"`        // Aggregate trade ID
	Price        float64 `json:"p,string"` // Price
	Quantity     float64 `json:"q,string"` // Quantity
	FirstTradeID int     `json:"f"`        // First trade ID
	LastTradeID  int     `json:"l"`        // Last trade ID
	TradeTime    int     `json:"T"`        // Trade time
	IsMaker      bool    `json:"m"`        // Is the buyer the market maker?
	Ignore       bool    `json:"M"`        // Ignore
}

type Client struct {
	Ctx context.Context
}

func NewClient(ctx context.Context) Client {
	return Client{Ctx: ctx}
}

func (c Client) wssBaseUrl() string {
	return "wss://stream.binance.com:9443/ws/"
}

func (c Client) ListenTrades(baseAsset string, quoteAsset string) (
	tradeCh chan *Trade,
	done chan struct{},
) {
	symbol := baseAsset + quoteAsset
	url := fmt.Sprintf(c.wssBaseUrl()+"%s@aggTrade", strings.ToLower(symbol))
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Binance websocket error:", err)
	}

	done = make(chan struct{})
	tradeCh = make(chan *Trade)

	go func() {
		defer func() {
			_ = conn.Close()
		}()
		defer closeChan(done)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Binance websocket reading error:", err)
				return
			}
			trade := &Trade{
				BaseAsset:  baseAsset,
				QuoteAsset: quoteAsset,
			}
			_ = json.Unmarshal(message, trade)
			tradeCh <- trade
		}
	}()

	go func() {
		defer func() {
			_ = conn.Close()
		}()
		defer closeChan(done)

		<-c.Ctx.Done()

		// Cleanly close the connection by sending a close message and then
		// waiting (with timeout) for the server to close the connection.
		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("Write close:", err)
		}
		log.Println("Closing websocket connection (Binance)")

		select {
		case <-done:
		case <-time.After(time.Second):
		}
	}()

	return
}

func closeChan(ch chan struct{}) {
	select {
	case <-ch:
		return
	default:
	}

	close(ch)
	return
}
