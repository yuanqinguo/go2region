package controls

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"go2region/utils"
	"go2region/utils/logs"
)

type Resp struct {
	Code   int         `json:"code"`
	Errmsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

//在这里统一封装回复
func WriterResp(ctx iris.Context, data interface{}, errno *utils.Errno, extMsg string) {
	formatExtMsg := func(extMsg string) string {
		if len(extMsg) > 0 {
			return fmt.Sprintf("[%s]", extMsg)
		}
		return ""
	}
	_, err := ctx.JSON(Resp{
		Code:   errno.Code,
		Errmsg: errno.Message + formatExtMsg(extMsg),
		Data:   data,
	})
	if err != nil {
		logs.LogSystem.Error("WriterResp error: ", err.Error())
	}
}
