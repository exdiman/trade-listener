package handler

import (
	"context"
	"log"
	"time"

	"github.com/exdiman/trade-listener/internal/app/model"
	"github.com/exdiman/trade-listener/internal/app/sqlstore"
)

func StoreTrades(ctx context.Context, tradeCh chan *model.Trade, store *sqlstore.Store) (done chan struct{}) {
	tradeRepo := store.TradeRepository()

	done = make(chan struct{})

	// обработка всех сделок со всех бирж из канала tradeCh
	go func() {
		defer close(done)

		bufMaxSize := int(cap(tradeCh) / 2)
		buf := make([]model.Trade, 0, bufMaxSize)
		sleep := make(chan struct{}, 1)
		defer close(sleep)

		saveBuf := func() {
			tradeRepo.Create(&buf)
			log.Printf("Insert %v trades\n", len(buf))
			buf = buf[:0]
		}

		for first := range tradeCh {
			buf = append(buf, *first)
		Loop:
			for len(buf) < bufMaxSize {
				select {
				case item := <-tradeCh:
					buf = append(buf, *item)
				case <-sleep:
					log.Println("Sleeping...")
					time.Sleep(2 * time.Second)
				case <-ctx.Done():
					log.Println("Inserting last trades...")
					saveBuf()
					return
				default:
					sleep <- struct{}{}
					break Loop
				}
			}

			saveBuf()
		}
	}()

	return
}
