package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/aisalamdag23/etherstats/internal/handler"
	ethhttp "github.com/aisalamdag23/etherstats/internal/handler/eth/v1"
	"github.com/aisalamdag23/etherstats/internal/infrastructure/config"
	"github.com/aisalamdag23/etherstats/internal/infrastructure/sql"
	"github.com/aisalamdag23/etherstats/internal/infrastructure/sql/postgres"
	ethdb "github.com/aisalamdag23/etherstats/internal/storage/db/eth"
	alchemysvc "github.com/aisalamdag23/etherstats/internal/usecase/alchemy"
	ethsvc "github.com/aisalamdag23/etherstats/internal/usecase/eth"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Registry is the factory that creates all the "feature servers"
type Registry struct {
	cfg     *config.Config
	db      *sqlx.DB
	redisDB *redis.Client
	logger  *zap.Logger
}

// Init instantiates the registry for API
// - creates database connection pool
func Init(ctx context.Context, cfg *config.Config, logger *zap.Logger) *Registry {
	registry := &Registry{
		cfg:    cfg,
		logger: logger,
	}

	// create a connection to db
	database, err := registry.createDB()
	if err != nil {
		logger.Fatal(err.Error())
	}

	registry.db = database

	// create a connection to redis
	redisDB, err := registry.createRedisDB(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}
	registry.redisDB = redisDB

	return registry
}

func (r *Registry) CreateETHServer() (handler.Handler, error) {
	repository := ethdb.NewRepository(r.db, r.redisDB, time.Second*time.Duration(r.cfg.Alchemy.CacheTTLSec))
	alchemySvc, err := alchemysvc.NewService(r.cfg.Alchemy.MainNetURL, r.cfg.Alchemy.APIKey)
	if err != nil {
		return nil, err
	}
	svc := ethsvc.NewService(repository, alchemySvc, r.logger)

	return ethhttp.NewServer(svc), nil
}

func (r *Registry) createDB() (*sqlx.DB, error) {
	dsnFactory := postgres.NewDSNFactory()
	dsn := dsnFactory.Create(
		r.cfg.PostgresDB.Credentials.Host,
		r.cfg.PostgresDB.Credentials.Port,
		r.cfg.PostgresDB.Credentials.User,
		r.cfg.PostgresDB.Credentials.Pass,
		r.cfg.PostgresDB.Credentials.DBName,
		r.cfg.PostgresDB.ConnectionTimeout,
	)

	connMaxLifetime, err := time.ParseDuration(fmt.Sprintf("%ds", r.cfg.PostgresDB.ConnLifetimeSec))
	if err != nil {
		return nil, err
	}

	dbFactory := sql.NewDBFactory()
	return dbFactory.Create(
		r.cfg.PostgresDB.Driver,
		dsn,
		r.cfg.PostgresDB.MaxOpenConn,
		connMaxLifetime,
	)
}

func (r *Registry) createRedisDB(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.cfg.RedisDB.Credentials.Host, r.cfg.RedisDB.Credentials.Port),
		Password: r.cfg.RedisDB.Credentials.Pass,
		DB:       0, // use default DB
		Protocol: 3,
	})

	// check if the connection is alive
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
