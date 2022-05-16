package api

import (
	"github.com/shapeshift/unchained/pkg/api"
	"github.com/shapeshift/unchained/pkg/cosmos"
)

// Contains info about the running coinstack
// swagger:model Info
type Info struct {
	// swagger:allOf
	api.BaseInfo
	// required: true
	// example: 291107634956378
	TotalSupply string `json:"totalSupply"`
	// required: true
	// example: 186039736185555
	BondedTokens string `json:"bondedTokens"`
	// required: true
	// example: 0.1541068456
	APR string `json:"apr"`
}

// Contains info about account details for an address or xpub
// swagger:model Account
type Account struct {
	// swagger:allOf
	api.BaseAccount
	// required: true
	// example: 420
	AccountNumber int `json:"accountNumber"`
	// required: true
	// example: 69
	Sequence int `json:"sequence"`
	// required: true
	Assets []cosmos.Value `json:"assets"`
}

// Contains info about a transaction
// swagger:model Tx
type Tx struct {
	// swagger:allOf
	api.BaseTx
	// required: true
	Confirmations int `json:"confirmations"`
	// required: true
	Fee cosmos.Value `json:"fee"`
	// required: true
	// example: 888
	GasUsed string `json:"gasUsed"`
	// required: true
	// example: 999
	GasWanted string `json:"gasWanted"`
	// required: true
	// example: 1
	Index int    `json:"index"`
	Memo  string `json:"memo,omitempty"`
	// required: true
	// 123456789
	Value string `json:"value"`
	// required: true
	Messages []cosmos.Message `json:"messages"`
	// required: true
	Events cosmos.EventsByMsgIndex `json:"events"`
}

// Contains info about transaction history for an address or xpub
// swagger:model TxHistory
type TxHistory struct {
	// swagger:allOf
	api.BaseTxHistory
	// required: true
	Txs []Tx `json:"txs"`
}
