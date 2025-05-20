package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

const configPathEnvName = "SPEC_FILE"

type (
	// Config ...
	Config struct {
		// CommitHash is a git commit hash of this app build
		CommitHash string
		// Tag is a git tag of this app build
		Tag        string
		General    General          `mapstructure:"general" validate:"required"`
		PostgresDB Database         `mapstructure:"postgresdb" validate:"required"`
		RedisDB    Database         `mapstructure:"redisdb" validate:"required"`
		Alchemy    APIProviderCreds `mapstructure:"alchemy" validate:"required"`
	}

	// General config.
	General struct {
		AppName string `mapstructure:"app_name" validate:"required"`
		// HTTPAddr internal server http address
		HTTPAddr        string `mapstructure:"http_addr" validate:"required"`
		WriteTimeoutSec int    `mapstructure:"http_write_timeout_sec" validate:"required"`
		ReadTimeoutSec  int    `mapstructure:"http_read_timeout_sec" validate:"required"`
		IdleTimeoutSec  int    `mapstructure:"http_idle_timeout_sec" validate:"required"`
		// ShutdownWaitSec is the number of secs the server will wait
		// before shutting down after it receives an exit signal
		ShutdownWaitSec int    `mapstructure:"graceful_shutdown_wait_time_sec" validate:"required"`
		LogLevel        string `mapstructure:"log_level" validate:"required"`
	}

	Database struct {
		Driver            string        `mapstructure:"driver"`
		Credentials       DBCredentials `mapstructure:"credentials" validate:"required"`
		ConnectionTimeout int           `mapstructure:"conn_timeout"`
		MaxOpenConn       int           `mapstructure:"max_open_conn"`
		ConnLifetimeSec   int           `mapstructure:"conn_lifetime_sec"`
	}

	DBCredentials struct {
		Host   string `mapstructure:"host" validate:"required"`
		Port   int    `mapstructure:"port" validate:"required"`
		DBName string `mapstructure:"name"`
		User   string `mapstructure:"user"`
		Pass   string `mapstructure:"pass"`
	}

	APIProviderCreds struct {
		APIKey      string `mapstructure:"api_key" validate:"required"`
		MainNetURL  string `mapstructure:"mainnet_url" validate:"required"`
		CacheTTLSec int    `mapstructure:"cache_ttl_sec" validate:"required"`
	}
)

// Load loads all configurations in to a new Config struct.
// CommitHash is a git commit hash of this app build.
// Tag is a git Tag of this app build.
func Load(commitHash, tag string) (*Config, error) {
	configFilePath := os.Getenv(configPathEnvName)
	if configFilePath == "" {
		return nil, fmt.Errorf("env variable %s is not defined", configPathEnvName)
	}

	v := viper.New()

	v.SetConfigFile(configFilePath)
	v.SetConfigType("yml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	c.CommitHash = commitHash
	c.Tag = tag

	validator := validator.New()
	err = validator.Struct(c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
