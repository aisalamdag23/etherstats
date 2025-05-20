package eth

import (
	"context"
	"time"

	"github.com/aisalamdag23/etherstats/internal/domain"
	"go.uber.org/zap"
)

type service struct {
	lgr            *zap.Logger
	repository     domain.Repository
	alchemyService domain.AlchemyAPIService
}

func NewService(repository domain.Repository, alchemyService domain.AlchemyAPIService, lgr *zap.Logger) domain.Service {
	return &service{
		repository:     repository,
		alchemyService: alchemyService,
		lgr:            lgr,
	}
}

// Get retrieves the gas price, latest block number, and balance of a given Ethereum address.
// It returns a Response struct containing the gas price, block number, balance, and server time.
func (s *service) Get(ctx context.Context, address string) (*domain.Response, error) {
	var response domain.Response
	// 1. Get the gas price
	gasPrice, err := s.getGasPrice(ctx)
	if err != nil {
		s.lgr.Error("failed to get gas price", zap.Error(err))
		return nil, err
	}
	response.GasPrice = gasPrice

	// 2. Get the latest block number
	blockNumber, err := s.getLatestBlockNumber(ctx)
	if err != nil {
		s.lgr.Error("failed to get latest block number", zap.Error(err))
		return nil, err
	}
	response.BlockNumber = blockNumber

	// 3. Get the balance of the address
	balance, err := s.getBalance(ctx, address)
	if err != nil {
		s.lgr.Error("failed to get balance", zap.Error(err), zap.String("address", address))
		return nil, err
	}
	response.Balance = *balance

	// 4. Set the server time
	response.ServerTime = time.Now().Format(time.RFC3339)

	return &response, nil
}

// getGasPrice retrieves the current gas price from the Ethereum network.
// It first checks if the gas price is cached in Redis.
// If not, it fetches the gas price from the Alchemy API and stores it in Redis.
// It returns the gas price as a string.
func (s *service) getGasPrice(ctx context.Context) (string, error) {
	// Check if the gas price is already cached in Redis
	// If not, fetch it from the Alchemy API and store it in Redis.
	price, err := s.repository.GetGasPrice(ctx)
	if err == nil && price != "" {
		// early return the cached gas price
		return price, nil
	}
	// If the gas price is not found in Redis, fetch it from the Alchemy API
	price, err = s.alchemyService.GetGasPrice(ctx)
	if err != nil {
		s.lgr.Error("failed to get gas price", zap.Error(err))
		return "", err
	}
	// Store the gas price in Redis with a TTL set from the config
	err = s.repository.SetGasPrice(ctx, price)
	if err != nil {
		s.lgr.Error("failed to set gas price in redis", zap.Error(err))
		// just log the error and return the price
	}

	return price, nil
}

// getLatestBlockNumber retrieves the latest block number from the Ethereum network.
// It first checks if the block number is cached in Redis.
// If not, it fetches the block number from the Alchemy API and stores it in Redis.
// It returns the block number as a uint64.
func (s *service) getLatestBlockNumber(ctx context.Context) (uint64, error) {
	// Check if the block number is already cached in Redis
	// If not, fetch it from the Alchemy API and store it in Redis.
	blockNumber, err := s.repository.GetBlockNumber(ctx)
	if err == nil && blockNumber != 0 {
		// early return the cached block number
		return blockNumber, nil
	}
	// If the block number is not found in Redis, fetch it from the Alchemy API
	blockNumber, err = s.alchemyService.GetLatestBlockNumber(ctx)
	if err != nil {
		s.lgr.Error("failed to get latest block number", zap.Error(err))
		return 0, err
	}
	// Store the block number in Redis with a TTL set from the config
	err = s.repository.SetBlockNumber(ctx, blockNumber)
	if err != nil {
		s.lgr.Error("failed to set block number in redis", zap.Error(err))
		// just log the error and return the block number
	}

	return blockNumber, nil
}

// getBalance retrieves the balance of a given Ethereum address in both Wei and Eth.
// It uses the Alchemy API service to fetch the balance and returns it as a Balance struct.
func (s *service) getBalance(ctx context.Context, address string) (*domain.Balance, error) {
	// Get the balance from the Alchemy API
	balance, err := s.alchemyService.GetBalance(ctx, address)
	if err != nil {
		s.lgr.Error("failed to get balance", zap.Error(err), zap.String("address", address))
		return nil, err
	}
	// Save the balance to the database
	_, err = s.repository.SaveBalance(ctx, address, balance)
	if err != nil {
		s.lgr.Error("failed to save balance", zap.Error(err), zap.String("address", address))
		// return the balance even if saving fails
	}

	return &domain.Balance{
		Address: address,
		Eth:     balance,
	}, nil
}
