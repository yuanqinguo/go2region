package utils

// 定义错误码
type Errno struct {
	Code    int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

/*
错误码设计
100 表示无错误
第一位表示错误级别, 1 为系统错误, 2 为数据库错误, 3xxx
第二位表示服务模块代码
第三位表示具体错误代码
*/

var (
	Ok = &Errno{Code: 100, Message: "OK"}

	// 系统错误, 前缀为 10
	ServerError = &Errno{Code: 101, Message: "内部服务器错误"}
	ErrArg      = &Errno{Code: 102, Message: "请求参数错误"}
)
