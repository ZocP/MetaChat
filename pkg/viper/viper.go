package viper

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"os"
)

const (
	configPath = "./files/config.yaml"
	configType = "yaml"
)

func NewViper() *viper.Viper {
	v := viper.New()

	v.SetConfigType(configType)
	v.SetConfigFile(configPath)
	_, err := os.Stat(configPath)
	if os.IsExist(err) {
		if _, err := os.Create(configPath); err != nil { // perm 0666
			panic(err)
		}
		if err := viper.SafeWriteConfig(); err != nil {
			panic(err)
		}
	}
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return v
}

func Provide() fx.Option {
	return fx.Provide(NewViper)
}
