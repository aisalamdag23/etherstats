package domain

import (
	"context"
	"time"
)

type (
	Service interface {
		Get(ctx context.Context, address string) (*Response, error)
	}

	Repository interface {
		SetGasPrice(ctx context.Context, price string) error
		SetBlockNumber(ctx context.Context, blockNumber uint64) error
		GetGasPrice(ctx context.Context) (string, error)
		GetBlockNumber(ctx context.Context) (uint64, error)
		SaveBalance(ctx context.Context, address, balance string) (*AddressBalance, error)
	}

	AlchemyAPIService interface {
		GetGasPrice(ctx context.Context) (string, error)
		GetLatestBlockNumber(ctx context.Context) (uint64, error)
		GetBalance(ctx context.Context, address string) (string, error)
	}

	Response struct {
		GasPrice    string  `json:"ethGasPrice"`
		BlockNumber uint64  `json:"latestBlockNumber"`
		Balance     Balance `json:"balance"`
		ServerTime  string  `json:"serverTime"`
	}

	Balance struct {
		Address string `json:"address"`
		Eth     string `json:"ethBalance"`
	}

	AddressBalance struct {
		ID        int       `db:"id"`
		Address   string    `db:"address"`
		Balance   string    `db:"balance"`
		CreatedAt time.Time `db:"created_at"`
	}
)
