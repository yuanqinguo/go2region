package controls

import (
	"github.com/kataras/iris/v12"
	"goip2region/utils"
	"goip2region/utils/helper"
)

type IpInfoResp struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Remark   string `json:"remark"`
}

func GetIpInfo(ctx iris.Context) {
	reqParams := ctx.URLParams()
	ipstr, _ := reqParams["ip"]
	err := utils.Ok
	var data interface{}
	if helper.CheckIp(ipstr) {
		ipdb := GetInstance()
		prv, city, region, remark := ipdb.GetIpInfo(helper.ConverIptoInt(ipstr))
		if len(prv) > 0 && len(city) > 0 {
			data = &IpInfoResp{
				Province: prv,
				City:     city,
				Region:   region,
				Remark:   remark,
			}
		}
	} else {
		err = utils.ErrArg
	}

	WriterResp(ctx, data, err, "")
}
