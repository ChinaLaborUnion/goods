package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/resource"
)

func RegisterResourceRouters(app *iris.Application) {

	// 资源路由
	resourceRouter := app.Party("goods/resources")

	resourceRouter.Get("/qiniu/upload_token", hero.Handler(resource.GetQiNiuUploadToken))

	//全局路由测试
	testRouter := app.Party("global/resources")
	testRouter.Get("/qiniu/upload_token", hero.Handler(resource.TestGetQiNiuUploadToken))
}
