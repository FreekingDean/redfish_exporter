package config

import "github.com/spf13/viper"

type ConfigOption func(*viper.Viper)

func WithFilePath(path string) ConfigOption {
	return func(v *viper.Viper) {
		v.SetConfigFile(path)
	}
}
