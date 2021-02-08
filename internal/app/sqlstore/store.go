package sqlstore

import (
	"github.com/jmoiron/sqlx"
	"log"
)

type Store struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) TradeRepository() *TradeRepository {
	return &TradeRepository{s}
}

func (s *Store) PrepareDb() {
	sql := `CREATE TABLE IF NOT EXISTS trades ( 
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT , 
			exchange ENUM('Binance','Poloniex') NOT NULL , 
			base_asset VARCHAR(4) NOT NULL , 
			quote_asset VARCHAR(4) NOT NULL , 
			trade_id VARCHAR(64) NOT NULL , 
			side ENUM('BUY','SELL') NOT NULL , 
			price DECIMAL(23,10) NOT NULL , 
			amount DECIMAL(23,10) NOT NULL , 
			created_at BIGINT NOT NULL , 
			PRIMARY KEY (id)
		) ENGINE = InnoDB DEFAULT CHARSET=utf8;`
	_, err := s.db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
}
