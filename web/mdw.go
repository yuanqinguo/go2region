package web

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"go2region/utils/logs"
	"runtime"
)

// 统一异常处理
func NewRecoverMdw() iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx.IsStopped() {
					return
				}

				var stacktrace string
				for i := 1; ; i++ {
					_, f, l, got := runtime.Caller(i)
					if !got {
						break
					}

					stacktrace += fmt.Sprintf("%s:%d\n", f, l)
				}

				// when stack finishes
				logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
				logMessage += fmt.Sprintf("At Request: %d %s %s %s\n", ctx.GetStatusCode(), ctx.Path(), ctx.Method(), ctx.RemoteAddr())
				logMessage += fmt.Sprintf("Trace: %s\n", err)
				logMessage += fmt.Sprintf("\n%s", stacktrace)

				logs.LogSystem.Errorf("recover => %s", logMessage)

				ctx.StatusCode(500)
				ctx.StopExecution()
			}
		}()

		ctx.Next()
	}
}
