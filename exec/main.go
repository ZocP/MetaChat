package main

import (
	"MetaChat/app"
	"MetaChat/app/metaChat"
	"MetaChat/pkg/http"
	"MetaChat/pkg/log"
	"MetaChat/pkg/router"
	"MetaChat/pkg/signal"
	"MetaChat/pkg/viper"
	"context"
	"go.uber.org/fx"
)

func main() {
	Launch()
}

func Launch() {
	//log := log.NewLogger(nil)
	fx.New(initPackages(), metaChat.Provide()).Run()
	//cmd := exec.Command("./files/cqhttp/go-cqhttp")
	//
	////cmd.Dir = "./files/cqhttp"
	//out, err := cmd.StdoutPipe()
	//if err != nil {
	//	log.Error("cq http not found, please init manually", zap.Error(err))
	//}
	//if err := cmd.Run(); err != nil {
	//	log.Error("cq http not found, please init manually", zap.Error(err))
	//}
	//for {
	//	tmp := make([]byte, 1024)
	//	o, err := out.Read(tmp)
	//	zap.S().Debug("output from cq: ", zap.String("info", string(rune(o))))
	//	if err != nil {
	//		break
	//	}
	//}
}

func initPackages() fx.Option {
	return fx.Options(
		http.Provide(),
		log.Provide(),
		signal.Provide(),
		viper.Provide(),
		router.Provide(),
		fx.Invoke(lc),
	)
}

func lc(lifecycle fx.Lifecycle, app app.APP, handler *signal.StopHandler) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return app.OnStart()
		},
		OnStop: func(ctx context.Context) error {
			handler.Stop()
			return nil
		},
	})
}
