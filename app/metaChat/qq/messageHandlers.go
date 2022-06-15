package qq

import (
	"MetaChat/pkg/qqBot"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type QQ struct {
	log *zap.Logger
}

func (qq *QQ) MessageHandler(ctx qqBot.Context, msg gjson.Result) {

}
