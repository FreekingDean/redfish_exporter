package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Hosts    map[string]Host `mapstructure:"hosts"`
	Groups   map[string]Host `mapstructure:"groups"`
	LogLevel string          `mapstructure:"logLevel"`
	Metrics  Metrics         `mapstructure:"metrics"`
	Web      Web             `mapstructure:"web"`
}

type Web struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
	Auth    Auth   `mapstructure:"auth"`
}

func (w Web) ListenAddress() string {
	return fmt.Sprintf("%s:%d", w.Address, w.Port)
}

type Auth struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Host struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`

	Metrics Metrics `mapstructure:"metrics"`
}

type Metrics struct {
	EnableAll bool              `mapstructure:"enableAll"`
	Metrics   map[string]Metric `mapstructure:"metrics"`
}

type Metric struct {
	Enabled bool              `mapstructure:"enabled"`
	Labels  map[string]string `mapstructure:"labels"`
}

func New(opts ...ConfigOption) (Config, error) {
	config := Config{}

	v := viper.New()
	v.SetDefault("logLevel", "info")
	v.SetDefault("web.address", "")
	v.SetDefault("web.port", 9100)
	v.SetDefault("web.auth.enabled", false)

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

	err := viper.Unmarshal(&config)

	return config, err
}
