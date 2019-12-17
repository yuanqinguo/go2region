package helper

import (
	"github.com/axgle/mahonia"
	"math/big"
	"net"
	"regexp"
	"strings"
)

func CheckIp(ip string) bool {
	address := net.ParseIP(ip)
	if address == nil {
		return false
	}
	return true
}

func ConverIptoInt(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func ZoomInTenThousand(ip int64) int64 {
	return ip * 10000
}

func DeleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)       //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}

func ConvertStr2GBK(str string) string {
	dec := mahonia.NewDecoder("GB18030")
	//converts a  string from UTF-8 to gbk encoding.
	return dec.ConvertString(str)
}

func ConvertGBK2Str(gbkStr string) string {
	enc := mahonia.NewEncoder("utf-8")
	//converts a  string from gbk to UTF-8 encoding.
	return enc.ConvertString(gbkStr)
}
