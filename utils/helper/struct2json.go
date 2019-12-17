package helper

import (
	"encoding/json"
	"goip2region/utils/logs"
)

func Struct2Json(st interface{}) (str string, err error) {
	b, e := json.Marshal(st)
	if e != nil {
		logs.LogSystem.Error("Struct2Json Error: ", err.Error())
	}
	str = string(b)
	e = err

	return
}
