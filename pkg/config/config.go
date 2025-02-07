package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Ethereum EthereumConfig
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type EthereumConfig struct {
	RPCURL        string        `mapstructure:"rpc_url"`
	RetryAttempts int           `mapstructure:"retry_attempts"`
	RetryDelay    time.Duration `mapstructure:"retry_delay"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("ethereum.rpc_url", "https://ethereum-rpc.publicnode.com")
	viper.SetDefault("ethereum.retry_attempts", 3)
	viper.SetDefault("ethereum.retry_delay", "2s")

	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ETH_PARSER")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
