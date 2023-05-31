package main

import "github.com/spf13/viper"

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	DB_SOURCE           string `mapstructure:"DB_SOURCE"`
	SUPABASE_JWT_SECRET string `mapstructure:"SUPABASE_JWT_SECRET"`
	REDIS_CONN_STRING   string `mapstructure:"REDIS_CONN_STRING"`
	// CLIENT_ID         string `mapstructure:"CLIENT_ID"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
