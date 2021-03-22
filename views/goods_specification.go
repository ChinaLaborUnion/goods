package views

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"grpc-demo/views/goods_specification"
)

func RegisterGoodsSpecificationRouters(app *iris.Application){
	templateRouter := app.Party("goods/specification_template")

	templateRouter.Post("", hero.Handler(goods_specification.CreateTemplate))
	templateRouter.Get("/list", hero.Handler(goods_specification.ListTemplate))
	templateRouter.Put("/{tid:int}", hero.Handler(goods_specification.PutTemplate))
	templateRouter.Delete("/{tid:int}", hero.Handler(goods_specification.DeleteTemplate))
	templateRouter.Post("/_mget", hero.Handler(goods_specification.MgetTemplate))
}
