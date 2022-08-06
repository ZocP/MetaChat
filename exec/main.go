package main

import (
	"MetaChat/app"
	"MetaChat/pkg/http"
	"MetaChat/pkg/log"
	"MetaChat/pkg/router"
	"MetaChat/pkg/signal"
	"MetaChat/pkg/viper"
	"context"
	"fmt"
	"go.uber.org/fx"
)

func main() {
	Launch()
}

func Launch() {
	//fx.New(initPackages(), metaChat.Provide()).Run()
	e5 := 6
	fmt.Println(e5)
	//qqBot.AddConditionHandler(condition.NewCondition(
	//	map[string]string{
	//		cq.SENDER_USERID: "1395437934",
	//	},
	//), func(ctx qqBot.Context, msg gjson.Result) {
	//	ctx.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetMessageQuick(msg, "test")))
	//})
	//fx.New(initPackages(), qqBot.Provide()).Run()
}

func initPackages() fx.Option {
	return fx.Options(
		http.Provide(),
		log.Provide(),
		signal.Provide(),
		viper.Provide(),
		router.Provide(),
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
