package viper

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	configPath = "./files/config.yaml"
	configType = "yaml"
)

func NewViper() *viper.Viper {
	v := viper.New()
	viper.GetViper()
	v.SetConfigType(configType)
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return v
}

func Provide() fx.Option {
	return fx.Provide(NewViper)
}
