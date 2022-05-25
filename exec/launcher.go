package main

import (
	"MetaChat/app/metaChat"
	"MetaChat/pkg/http"
	"MetaChat/pkg/log"
	"MetaChat/pkg/signal"
	"MetaChat/pkg/viper"
	"context"
	"go.uber.org/fx"
)

func Launch() {
	fx.New(initPackages(), metaChat.Provide()).Run()
}

func initPackages() fx.Option {
	return fx.Options(
		http.Provide(),
		log.Provide(),
		signal.Provide(),
		viper.Provide(),
		fx.Invoke(lc),
	)
}

func lc(lifecycle fx.Lifecycle, server *http.Server, app *metaChat.MetaChat, handler *signal.StopHandler) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := app.OnStart(); err != nil {
				return err
			}
			return server.Start()
		},
		OnStop: func(ctx context.Context) error {
			handler.Stop()
			return nil
		},
	})
}
