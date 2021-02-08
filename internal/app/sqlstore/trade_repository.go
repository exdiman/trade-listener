package sqlstore

import (
	"log"

	"github.com/exdiman/trade-listener/internal/app/model"
)

type TradeRepository struct {
	store *Store
}

func (r TradeRepository) Create(trades *[]model.Trade) {
	sql := `INSERT INTO trades 
		(exchange, base_asset, quote_asset, trade_id, side, price, amount, created_at) 
		VALUES (:exchange, :base_asset, :quote_asset, :trade_id, :side, :price, :amount, :created_at)`

	_, err := r.store.db.NamedExec(sql, *trades)
	if err != nil {
		log.Fatal(err)
	}
}
