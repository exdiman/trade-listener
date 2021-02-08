package model

type Exchange string

type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

type Trade struct {
	Exchange    Exchange `db:"exchange"`
	BaseAsset   string   `db:"base_asset"`
	QuoteAsset  string   `db:"quote_asset"`
	TradeID     string   `db:"trade_id"`
	Side        Side     `db:"side"`
	Price       float64  `db:"price"`
	Amount      float64  `db:"amount"`
	TimestampMs int      `db:"created_at"`
}
