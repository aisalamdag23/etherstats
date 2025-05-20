package alchemy

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/aisalamdag23/etherstats/internal/domain"

	// Importing the go-ethereum package for Ethereum client
	// and common types
	// This package provides the necessary functions to interact with the Ethereum blockchain
	// and perform operations like fetching gas prices, block numbers, and balances.
	// Visit: https://geth.ethereum.org/docs/developers/dapp-developer/native
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type service struct {
	client *ethclient.Client
}

func NewService(baseURL, apiKey string) (domain.AlchemyAPIService, error) {
	client, err := ethclient.Dial(baseURL + "/" + apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	return &service{
		client: client,
	}, nil
}

// GetGasPrice fetches the current suggested gas price from the Ethereum network.
// It returns the gas price in ETH for human readability
func (s *service) GetGasPrice(ctx context.Context) (string, error) {
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch gas price: %v", err)
	}

	return s.convertToETH(gasPrice), nil
}

// GetLatestBlockNumber fetches the latest block number from the Ethereum network.
// It returns the block number as a uint64.
func (s *service) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := s.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch block number: %v", err)
	}

	return blockNumber, nil
}

// GetBalance fetches the balance of a given Ethereum address.
// It returns the balance in ETH for human readability
func (s *service) GetBalance(ctx context.Context, address string) (string, error) {
	addr := common.HexToAddress(address)
	balanceWei, err := s.client.BalanceAt(ctx, addr, nil) // nil = latest block

	return s.convertToETH(balanceWei), err
}

func (s *service) convertToETH(val *big.Int) string {
	flt := new(big.Float).SetInt(val)
	// Convert wei → ETH for human readability
	eth := new(big.Float).Quo(flt, big.NewFloat(math.Pow10(18)))

	return fmt.Sprintf("%.18f", eth)
}
