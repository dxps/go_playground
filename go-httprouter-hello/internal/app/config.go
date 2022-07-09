package app

import "github.com/spf13/viper"

// LocalConfig is a read-only configuration,
// whose source is either a file, environment variables, or both.
type LocalConfig struct {
	API_host string `mapstructure:"API_HTTP_HOST"` // The hostname it listens for HTTP API requests.
	API_port int    `mapstructure:"API_HTTP_PORT"` // The port number it listens for HTTP API requests.
}

func NewConfig() (*LocalConfig, error) {

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	c := &LocalConfig{}
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}
	return c, nil
}
