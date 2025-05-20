package eth

import (
	"context"
	"strconv"
	"time"

	"github.com/aisalamdag23/etherstats/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type (
	repository struct {
		db       *sqlx.DB
		redisDB  *redis.Client
		cacheTTL time.Duration
	}
)

const (
	// gasPriceKey is the key used to store the gas price in Redis
	gasPriceKey = "gas_price"
	// blockNumberKey is the key used to store the block number in Redis
	blockNumberKey = "block_number"
)

func NewRepository(db *sqlx.DB, redisDB *redis.Client, cacheTTL time.Duration) domain.Repository {
	return &repository{
		db:       db,
		redisDB:  redisDB,
		cacheTTL: cacheTTL,
	}
}

// SetGasPrice sets the current gas price in Redis with a specified TTL.
// It returns an error if the operation fails.
func (r *repository) SetGasPrice(ctx context.Context, price string) error {
	return r.redisDB.Set(ctx, gasPriceKey, price, r.cacheTTL).Err()
}

// SetBlockNumber sets the latest block number in Redis with a specified TTL.
// It returns an error if the operation fails.
func (r *repository) SetBlockNumber(ctx context.Context, blockNumber uint64) error {
	return r.redisDB.Set(ctx, blockNumberKey, blockNumber, r.cacheTTL).Err()
}

// GetGasPrice retrieves the current gas price from Redis.
// If the value is not found, it returns an empty string.
func (r *repository) GetGasPrice(ctx context.Context) (string, error) {
	val, err := r.redisDB.Get(ctx, gasPriceKey).Result()
	if err != nil {
		if err != redis.Nil {
			return "", err
		}

		// If the value is not found in Redis, return an empty string
		return "", nil
	}

	return val, nil
}

// GetBlockNumber retrieves the latest block number from Redis.
// If the value is not found, it returns 0.
func (r *repository) GetBlockNumber(ctx context.Context) (uint64, error) {
	val, err := r.redisDB.Get(ctx, blockNumberKey).Result()
	if err != nil {
		if err != redis.Nil {
			return 0, err
		}

		// If the value is not found in Redis, return 0
		return 0, nil
	}

	blockNumber, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

// SaveBalance saves the balance of an Ethereum address to the database.
// It returns the saved AddressBalance object or an error if the operation fails.
func (r *repository) SaveBalance(ctx context.Context, address, balance string) (*domain.AddressBalance, error) {
	query := `INSERT INTO balances 
				(address, balance)
			  VALUES 
				($1, $2)
			  RETURNING *;`

	var bal domain.AddressBalance
	err := r.db.GetContext(ctx, &bal, query, address, balance)
	if err != nil {
		return nil, err
	}

	return &bal, nil
}
