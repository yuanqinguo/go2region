package routes

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"goip2region/controls"
)

// 定义500错误处理函数
func err500(ctx iris.Context) {
	_, _ = ctx.WriteString("CUSTOM 500 ERROR")
}

// 定义404错误处理函数
func err404(ctx iris.Context) {
	_, _ = ctx.WriteString("CUSTOM 404 ERROR")
}

// 跨域处理？主域名统一，是否有必要
func corsHandler(app *iris.Application) {
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, //允许通过的主机名称,全部或者白名单
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		//AllowCredentials: true,
	})
	app.Use(crs)
	app.AllowMethods(iris.MethodOptions)
}

// 注入路由
func InnerRoute(app *iris.Application) {
	app.OnErrorCode(iris.StatusInternalServerError, err500)
	app.OnErrorCode(iris.StatusNotFound, err404)

	root := app.Party("/goip2region")
	root.Get("/ping", func(ctx iris.Context) { _, _ = ctx.WriteString("pong") })
	root.Get("/ipinfo", controls.GetIpInfo)
	root.Post("/reload", controls.Reloader)
}
