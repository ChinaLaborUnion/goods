package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/core/cache"
	viewbase "grpc-demo/core/view"
	"grpc-demo/models/db"
	"grpc-demo/utils"
	"grpc-demo/utils/middlewares"
	"grpc-demo/views"
)

func initRouter(app *iris.Application) {
	views.RegisterGoodsTagRouters(app)
	views.RegisterGoodsSpecificationRouters(app)
	views.RegisterGoodsRouters(app)
	views.RegisterResourceRouters(app)
	views.RegisterGoodsSlideshowRouters(app)
	views.RegisterGlobalRouters(app)
	views.RegisterOrderRouters(app)
	views.RegisterTestRouters(app)
	views.RegisterCourseRouters(app)
	views.RegisterCourseTagRouters(app)
	views.RegisterCourseApplyRouters(app)
	views.RegisterMsetRouters(app)
}

func main() {
	app := iris.New()
	// 注册控制器
	app.UseGlobal(middlewares.AbnormalHandle, middlewares.RequestLogHandle)
	hero.Register(viewbase.ViewBase)
	// 注册路由
	initRouter(app)
	// 初始化配置
	utils.InitGlobal()
	// 初始化数据库
	db.InitDB()
	// 初始化缓存
	//cache.InitDijan()
	cache.InitRedisPool()
	// 初始化任务队列
	//queue.InitTaskQueue()
	// 启动系统
	app.Run(iris.Addr(":80"), iris.WithoutServerError(iris.ErrServerClosed))
}
