package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Host     Host    `mapstructure:"host"`
	LogLevel string  `mapstructure:"logLevel"`
	Metrics  Metrics `mapstructure:"metrics"`
	Web      Web     `mapstructure:"web"`
}

type Web struct {
	Address    string `mapstructure:"address"`
	Port       int    `mapstructure:"port"`
	ConfigFile string `mapstructure:"configFile"`
}

func (w Web) ListenAddress() string {
	return fmt.Sprintf("%s:%d", w.Address, w.Port)
}

type Host struct {
	Endpoint string `mapstructure:"endpoint"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Metrics struct {
	EnableAll bool              `mapstructure:"enableAll"`
	Metrics   map[string]Metric `mapstructure:"metrics"`
}

type Metric struct {
	Enabled bool              `mapstructure:"enabled"`
	Labels  map[string]string `mapstructure:"labels"`
}

func New(opts ...Option) (Config, error) {
	config := Config{}

	v := viper.New()
	v.SetDefault("logLevel", "info")
	v.SetDefault("web.address", "")
	v.SetDefault("web.port", 9610)

	v.SetConfigType("yaml")
	v.SetConfigFile("./config.yaml")

	v.SetEnvPrefix("REDFISH_EXPORTER")

	for _, opt := range opts {
		opt(v)
	}

	if err := v.ReadInConfig(); err != nil {
		return config, err
	}
	v.AutomaticEnv()

	err := v.Unmarshal(&config)

	return config, err
}
