package config

import "github.com/spf13/viper"

type Config struct {
	Mode string `json:"mode"`
}

func Unmarshal(v *viper.Viper) *Config {
	return &Config{}
}
