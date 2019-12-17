package controls

import (
	"github.com/kataras/iris/v12"
	. "goip2region/config"
	"goip2region/utils"
)

func Reloader(ctx iris.Context) {
	baseConf := Config.BaseConf

	go GetInstance().Reloader(baseConf.IpdataPath)

	WriterResp(ctx, "", utils.Ok, "")
}
