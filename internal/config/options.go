package config

import "github.com/spf13/viper"

type Option func(*viper.Viper)

func WithFilePath(path string) Option {
	return func(v *viper.Viper) {
		v.SetConfigFile(path)
	}
}
