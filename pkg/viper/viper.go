package viper

import "github.com/spf13/viper"

func NewViper() *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile("./files/config.yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return v
}
