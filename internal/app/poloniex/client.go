package poloniex

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Trade struct {
	BaseAsset  string
	QuoteAsset string
	TradeID    string
	Side       float64
	Price      string
	Size       string
	Timestamp  float64
}

type Client struct {
	Ctx context.Context
}

func (c Client) wssBaseUrl() string {
	return "wss://api2.poloniex.com"
}

func (c Client) sendMessage(conn *websocket.Conn, message []byte) {
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Write message error:", err)
	}
}

func (c Client) subscribeToSymbol(conn *websocket.Conn, symbol string) {
	message := []byte(fmt.Sprintf(`{"command": "subscribe", "channel": "%s"}`, symbol))
	c.sendMessage(conn, message)
}

func (c Client) ListenTrades(baseAsset string, quoteAsset string) (
	tradeCh chan *Trade,
	done chan struct{},
) {
	symbol := baseAsset + "_" + quoteAsset
	conn, _, err := websocket.DefaultDialer.Dial(c.wssBaseUrl(), nil)
	if err != nil {
		log.Fatal("Poloniex websocket error:", err)
	}

	c.subscribeToSymbol(conn, symbol)

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
				log.Println("Poloniex websocket reading error:", err)
				return
			}

			var data []interface{}
			_ = json.Unmarshal(message, &data)

			if len(data) < 3 {
				continue
			}

			switch data[2].(type) {
			case []interface{}:
				events := data[2].([]interface{})
				for _, v := range events {
					event := v.([]interface{})
					switch event[0].(type) {
					case string:
						if event[0] == "t" {
							//["t", "<trade id>", <1 for buy 0 for sell>, "<price>", "<size>", <timestamp>]
							trade := &Trade{
								baseAsset,
								quoteAsset,
								event[1].(string),
								event[2].(float64),
								event[3].(string),
								event[4].(string),
								event[5].(float64),
							}
							tradeCh <- trade
						}
					}
				}
			}
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
		log.Println("Closing websocket connection (Poloniex)")

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
