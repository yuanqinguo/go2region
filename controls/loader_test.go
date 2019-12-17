package controls

import (
	"fmt"
	"strings"
	"testing"
)

func TestHandleValue(t *testing.T) {
	mask := "台湾省嘉义县"
	prov := ""
	city := ""
	region := ""
	index := strings.Index(mask, "省")
	if index > 0 {
		prov = mask[:index+3]
		cityIndex := strings.Index(mask, "市")
		if cityIndex >= 0 {
			city = mask[index+3 : cityIndex+3]
			region = mask[cityIndex+3:]
		} else {
			cityIndex := strings.LastIndex(mask, "州")
			if cityIndex >= index+3 {
				city = mask[index+3 : cityIndex+3]
				region = mask[cityIndex+3:]
			} else {
				cityIndex := strings.Index(mask, "地区")
				if cityIndex >= index+3 {
					city = mask[index+3 : cityIndex+6]
					region = mask[cityIndex+6:]
				} else {
					// 有些数据直接省后面接了区县
					regionIndex := strings.Index(mask, "县")
					if regionIndex >= index+3 {
						region = mask[index+3 : regionIndex+3]
					}
				}
			}
		}
	}
	fmt.Printf("p: %s, c: %s, r: %s", prov, city, region)
}
