package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/global"
)

func RegisterGlobalRouters(app *iris.Application){
	//icon
	iconRouter := app.Party("goods/icon")

	iconRouter.Post("", hero.Handler(global.Create))
	iconRouter.Get("/list", hero.Handler(global.List))
	iconRouter.Put("/{id:int}", hero.Handler(global.Put))
	iconRouter.Delete("/{id:int}", hero.Handler(global.Delete))
	iconRouter.Post("/_mget", hero.Handler(global.Mget))

	setTimeRouter := app.Party("global/set_time")
	setTimeRouter.Post("", hero.Handler(global.SetTime))
	setTimeRouter.Get("", hero.Handler(global.GetTime))

}
