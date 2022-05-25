package config

import "github.com/spf13/viper"

type Config struct {
}

func Unmarshal(v *viper.Viper) *Config {
	return &Config{}
}
